package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/repository"
	"github.com/pastorenue/kinance/internal/user"
	"github.com/pastorenue/kinance/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type GoogleHandler struct {
	conf      *oauth2.Config
	repo      *repository.TokenRepository
	userSvc   *user.Service
	authSvc   *Service
	allowedHD string
}

func NewGoogleHandler(cfg *config.Config, repo *repository.TokenRepository, userSvc *user.Service, authSvc *Service) *GoogleHandler {
	conf := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return &GoogleHandler{conf: conf, repo: repo, userSvc: userSvc, authSvc: authSvc, allowedHD: cfg.Google.AllowedHD}
}

// GET /api/v1/auth/google/login
func (h *GoogleHandler) Login(c *gin.Context) {
	if h.conf.ClientID == "" || h.conf.RedirectURL == "" {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Google OAuth misconfigured: set GOOGLE_CLIENT_ID and GOOGLE_REDIRECT_URL",
		})
		return
	}

	state := generateRandomString(32)
	_ = h.repo.StoreState(state, time.Now().Add(10*time.Minute))
	url := h.conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("redirect_uri", h.conf.RedirectURL), // be explicit to avoid mismatch
		oauth2.SetAuthURLParam("prompt", "select_account consent"),
	)
	c.Redirect(http.StatusFound, url)
}

// GET /api/v1/auth/google/callback
func (h *GoogleHandler) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		c.JSON(http.StatusBadRequest, common.APIResponse{Success: false, Error: "missing state or code"})
		return
	}
	if ok := h.repo.ConsumeState(state); !ok {
		c.JSON(http.StatusBadRequest, common.APIResponse{Success: false, Error: "invalid state"})
		return
	}

	tok, err := h.conf.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{Success: false, Error: "failed to exchange code"})
		return
	}

	rawID, ok := tok.Extra("id_token").(string)
	if !ok || rawID == "" {
		c.JSON(http.StatusBadRequest, common.APIResponse{Success: false, Error: "id_token not present in token response"})
		return
	}

	payload, err := idtoken.Validate(c.Request.Context(), rawID, h.conf.ClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.APIResponse{Success: false, Error: "invalid id_token"})
		return
	}

	claims := payload.Claims
	email, _ := claims["email"].(string)
	emailVerified, _ := claims["email_verified"].(bool)
	given, _ := claims["given_name"].(string)
	family, _ := claims["family_name"].(string)
	if email == "" || !emailVerified {
		c.JSON(http.StatusUnauthorized, common.APIResponse{Success: false, Error: "email not verified"})
		return
	}
	if h.allowedHD != "" {
		if hd, _ := claims["hd"].(string); hd != h.allowedHD {
			c.JSON(http.StatusUnauthorized, common.APIResponse{Success: false, Error: "unauthorized domain"})
			return
		}
	}

	// Upsert user
	u, err := h.userSvc.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			req := &user.CreateUserRequest{Email: email, FirstName: given, LastName: family}
			u, err = h.userSvc.CreateUserWithoutPassword(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.APIResponse{Success: false, Error: err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, common.APIResponse{Success: false, Error: err.Error()})
			return
		}
	}

	access, err := h.authSvc.generateAccessToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{Success: false, Error: "failed to generate access token"})
		return
	}
	refresh, err := h.authSvc.generateRefreshToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{Success: false, Error: "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"access_token":  access,
			"refresh_token": refresh,
			"expires_in":    h.authSvc.jwtConfig.ExpirationTime,
			"user":          u,
		},
	})
}

// generateRandomString returns a URL-safe random string of exact length using crypto/rand.
// It uses base64.RawURLEncoding and sizes the random byte slice to guarantee enough output.
func generateRandomString(length int) string {
	if length <= 0 {
		return ""
	}
	// bytesNeeded so that encoded string (without padding) has at least length chars
	bytesNeeded := (length*3 + 3) / 4 // ceil(length*3/4)
	b := make([]byte, bytesNeeded)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	s := base64.RawURLEncoding.EncodeToString(b)
	if len(s) < length {
		// generate more and append (extremely unlikely)
		extra := make([]byte, bytesNeeded)
		if _, err := rand.Read(extra); err == nil {
			s += base64.RawURLEncoding.EncodeToString(extra)
		}
	}
	return s[:length]
}
