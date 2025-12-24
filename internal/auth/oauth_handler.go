package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/model"
	"github.com/pastorenue/kinance/internal/user"
)

type OAuthHandler struct {
	svc     *OAuthService
	userSvc *user.Service
}

func NewOAuthHandler(svc *OAuthService, userSvc *user.Service) *OAuthHandler {
	return &OAuthHandler{svc: svc, userSvc: userSvc}
}

// Optional convenience: register user via JSON (kept for parity with your previous code)
func (h *OAuthHandler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{Success: false, Error: err.Error()})
		return
	}
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, common.APIResponse{Success: false, Error: fmt.Errorf("invalid data").Error()})
		return
	}
	_, err := h.userSvc.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, common.APIResponse{Success: true, Message: "User created"})
}

// GET /api/v1/oauth/authorize (gin.WrapF)
func (h *OAuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	responseType := r.FormValue("response_type")
	scope := r.FormValue("scope")
	state := r.FormValue("state")

	if responseType != "code" {
		h.errorResponse(w, "unsupported_response_type", http.StatusBadRequest)
		return
	}

	// TODO: replace with actual authenticated user id
	userID := "user123"

	code, err := h.svc.GenerateAuthorizationCode(clientID, redirectURI, scope, userID)
	if err != nil {
		h.errorResponse(w, "server_error", http.StatusInternalServerError)
		return
	}

	redirectURL := redirectURI + "?code=" + code
	if state != "" {
		redirectURL += "&state=" + state
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// POST /api/v1/oauth/token (gin.WrapF)
func (h *OAuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	grantType := r.FormValue("grant_type")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	if err := h.svc.ValidateClient(clientID, clientSecret); err != nil {
		h.errorResponse(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	var token *model.Token
	var err error
	switch model.GrantType(grantType) {
	case model.GrantTypeAuthorisationCode:
		code := r.FormValue("code")
		redirectURI := r.FormValue("redirect_uri")
		token, err = h.svc.ExchangeCode(code, clientID, redirectURI)
	case model.GrantTypeRefreshToken:
		refreshToken := r.FormValue("refresh_token")
		token, err = h.svc.RefreshAccessToken(refreshToken, clientID)
	default:
		h.errorResponse(w, "unsupported_grant_type", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.errorResponse(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	h.jsonResponse(w, token, http.StatusOK)
}

// POST /api/v1/oauth/introspect (gin.WrapF)
func (h *OAuthHandler) Introspect(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	resp, err := h.svc.IntrospectToken(token)
	if err != nil {
		h.errorResponse(w, "server_error", http.StatusInternalServerError)
		return
	}
	h.jsonResponse(w, resp, http.StatusOK)
}

// POST /api/v1/oauth/revoke (gin.WrapF)
func (h *OAuthHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")
	if err := h.svc.ValidateClient(clientID, clientSecret); err != nil {
		h.errorResponse(w, "invalid_client", http.StatusUnauthorized)
		return
	}
	_ = h.svc.RevokeToken(token)
	w.WriteHeader(http.StatusOK)
}

func (h *OAuthHandler) jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *OAuthHandler) errorResponse(w http.ResponseWriter, error string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": error})
}
