package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pastorenue/kinance/internal/model"
	"github.com/pastorenue/kinance/internal/repository"
	"github.com/pastorenue/kinance/pkg/config"
)

var (
	ErrInvalidClient       = errors.New("invalid client")
	ErrInvalidGrant        = errors.New("invalid grant")
	ErrInvalidCode         = errors.New("invalid authorization code")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type OAuthService struct {
	cfg  *config.Config
	repo *repository.TokenRepository
}

func NewOAuthService(cfg *config.Config, repo *repository.TokenRepository) *OAuthService {
	return &OAuthService{cfg: cfg, repo: repo}
}

func (s *OAuthService) ValidateClient(clientID, clientSecret string) error {
	if clientID != s.cfg.Client.ClientID || clientSecret != s.cfg.Client.ClientSecret {
		return ErrInvalidClient
	}
	return nil
}

func (s *OAuthService) GenerateAuthorizationCode(clientID, redirectURI, scope, userID string) (string, error) {
	code := generateRandomString(32)
	authCode := &model.AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		UserID:      userID,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
	if err := s.repo.StoreCode(authCode); err != nil {
		return "", err
	}
	return code, nil
}

func (s *OAuthService) ExchangeCode(code, clientID, redirectURI string) (*model.Token, error) {
	authCode, err := s.repo.GetCode(code)
	if err != nil {
		return nil, ErrInvalidCode
	}
	if authCode.ClientID != clientID || authCode.RedirectURI != redirectURI {
		return nil, ErrInvalidGrant
	}
	_ = s.repo.DeleteCode(code)
	return s.generateTokens(clientID, authCode.UserID, authCode.Scope)
}

func (s *OAuthService) RefreshAccessToken(refreshToken, clientID string) (*model.Token, error) {
	stored, err := s.repo.GetToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	if stored.ClientID != clientID {
		return nil, ErrInvalidClient
	}
	_ = s.repo.RevokeToken(refreshToken)
	return s.generateTokens(clientID, stored.UserID, stored.Scope)
}

func (s *OAuthService) IntrospectToken(token string) (*model.IntrospectionResponse, error) {
	stored, err := s.repo.GetToken(token)
	if err != nil {
		return &model.IntrospectionResponse{Active: false}, nil
	}
	return &model.IntrospectionResponse{
		Active:    true,
		Scope:     stored.Scope,
		ClientID:  stored.ClientID,
		UserID:    stored.UserID,
		TokenType: string(model.TokenTypeBearer),
		Exp:       stored.ExpiresAt.Unix(),
	}, nil
}

func (s *OAuthService) RevokeToken(token string) error {
	return s.repo.RevokeToken(token)
}

func (s *OAuthService) generateTokens(clientID, userID, scope string) (*model.Token, error) {
	accessToken, err := s.generateJWT(clientID, userID, scope, time.Second*time.Duration(s.cfg.JWT.ExpirationTime))
	if err != nil {
		return nil, err
	}
	refreshToken := generateRandomString(64)
	stored := &model.StoredToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClientID:     clientID,
		UserID:       userID,
		Scope:        scope,
		ExpiresAt:    time.Now().Add(time.Second * time.Duration(s.cfg.JWT.ExpirationTime)),
		CreatedAt:    time.Now(),
	}
	if err := s.repo.StoreToken(stored); err != nil {
		return nil, err
	}
	return &model.Token{
		AccessToken:  accessToken,
		TokenType:    model.TokenTypeBearer,
		ExpiresIn:    int64(s.cfg.JWT.ExpirationTime),
		RefreshToken: refreshToken,
		Scope:        scope,
	}, nil
}

func (s *OAuthService) generateJWT(clientID, userID, scope string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"client_id": clientID,
		"user_id":   userID,
		"scope":     scope,
		"exp":       time.Now().Add(duration).Unix(),
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}
