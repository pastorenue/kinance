package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pastorenue/kinance/internal/model"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrCodeNotFound  = errors.New("authorization code not found")
	ErrInvalidCode   = errors.New("authorization code cannot be empty string")
)

type TokenRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewTokenRepository(redisAddr, redisPassword string, redisDB int) *TokenRepository {
	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword,
		DB:           redisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})
	return &TokenRepository{
		client: client,
		ctx:    context.Background(),
	}
}
func (r *TokenRepository) tokenKey(token string) string {
	return fmt.Sprintf("oauth:token:%s", token)
}

func (r *TokenRepository) codeKey(code string) string {
	return fmt.Sprintf("oauth:code:%s", code)
}

func (r *TokenRepository) stateKey(state string) string {
	return fmt.Sprintf("oauth:state:%s", state)
}

func (r *TokenRepository) Close() error {
	return r.client.Close()
}

func (r *TokenRepository) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

func (r *TokenRepository) StoreToken(token *model.StoredToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return ErrTokenNotFound
	}

	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return ErrTokenNotFound
	}

	pipe := r.client.Pipeline()

	// Store the access Token
	pipe.Set(r.ctx, r.tokenKey(token.AccessToken), data, ttl)

	// Store refresh token with reference
	if token.RefreshToken != "" {
		pipe.Set(r.ctx, r.tokenKey(token.RefreshToken), data, ttl)
	}
	_, err = pipe.Exec(r.ctx)
	return err
}

func (r *TokenRepository) GetToken(token string) (*model.StoredToken, error) {
	data, err := r.client.Get(r.ctx, r.tokenKey(token)).Bytes()
	if err == redis.Nil {
		return nil, ErrTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var stored model.StoredToken
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	if time.Now().After(stored.ExpiresAt) {
		r.client.Del(r.ctx, r.tokenKey(token))
		return nil, ErrTokenNotFound
	}

	return &stored, nil
}

func (r *TokenRepository) RevokeToken(token string) error {
	stored, err := r.GetToken(token)
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()
	pipe.Del(r.ctx, r.tokenKey(stored.AccessToken))
	if stored.RefreshToken != "" {
		pipe.Del(r.ctx, r.tokenKey(stored.RefreshToken))
	}

	_, err = pipe.Exec(r.ctx)
	return err
}

func (r *TokenRepository) StoreCode(code *model.AuthorizationCode) error {
	data, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("failed to marshal code: %w", err)
	}

	ttl := time.Until(code.ExpiresAt)
	if ttl <= 0 {
		return errors.New("code already expired")
	}

	return r.client.Set(r.ctx, r.codeKey(code.Code), data, ttl).Err()
}

func (r *TokenRepository) GetCode(code string) (*model.AuthorizationCode, error) {
	data, err := r.client.Get(r.ctx, r.codeKey(code)).Bytes()
	if err == redis.Nil {
		return nil, ErrCodeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get code: %w", err)
	}

	var authCode model.AuthorizationCode
	if err := json.Unmarshal(data, &authCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code: %w", err)
	}

	if time.Now().After(authCode.ExpiresAt) {
		r.client.Del(r.ctx, r.codeKey(code))
		return nil, ErrCodeNotFound
	}

	return &authCode, nil
}

func (r *TokenRepository) DeleteCode(code string) error {
	return r.client.Del(r.ctx, r.codeKey(code)).Err()
}

// StoreState stores a transient OAuth2 state with TTL
func (r *TokenRepository) StoreState(state string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("state already expired")
	}
	return r.client.Set(r.ctx, r.stateKey(state), "1", ttl).Err()
}

// ConsumeState deletes the state key if present and returns true if it existed
func (r *TokenRepository) ConsumeState(state string) bool {
	res := r.client.Del(r.ctx, r.stateKey(state))
	if err := res.Err(); err != nil {
		return false
	}
	return res.Val() > 0
}
