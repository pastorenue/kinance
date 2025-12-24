package model

import "time"

type GrantType string

const (
	GrantTypeAuthorisationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
	GrantTypeClientCredentials GrantType = "client_credentials"
)

type TokenType string

const (
	TokenTypeBearer TokenType = "Bearer"
)

type AuthorizationCode struct {
	Code        string
	ClientID    string
	RedirectURI string
	Scope       string
	UserID      string
	ExpiresAt   time.Time
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    TokenType `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Scope        string    `json:"scope,omitempty"`
}

type StoredToken struct {
	AccessToken  string
	RefreshToken string
	ClientID     string
	UserID       string
	Scope        string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type IntrospectionResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
}
