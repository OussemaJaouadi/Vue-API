package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type TokenConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	Issuer        string
}

type TokenClaims struct {
	UserID       string `json:"sub"`
	Type         string `json:"typ"`
	Audience     string `json:"aud"`
	Issuer       string `json:"iss"`
	ExpiresAt    int64  `json:"exp"`
	IssuedAt     int64  `json:"iat"`
	NotBefore    int64  `json:"nbf"`
	TokenVersion int    `json:"token_version"`
}

type TokenManager struct {
	config TokenConfig
	now    func() time.Time
}

func NewTokenManager(config TokenConfig) TokenManager {
	return TokenManager{
		config: config,
		now:    time.Now,
	}
}

func (manager TokenManager) IssueAccessToken(userID string, tokenVersion int) (string, error) {
	return manager.issueToken(userID, tokenVersion, TokenTypeAccess, "api", manager.config.AccessTTL, []byte(manager.config.AccessSecret))
}

func (manager TokenManager) IssueRefreshToken(userID string, tokenVersion int) (string, error) {
	return manager.issueToken(userID, tokenVersion, TokenTypeRefresh, "auth-refresh", manager.config.RefreshTTL, []byte(manager.config.RefreshSecret))
}

func (manager TokenManager) ValidateAccessToken(token string) (TokenClaims, error) {
	return manager.validateToken(token, TokenTypeAccess, "api", []byte(manager.config.AccessSecret))
}

func (manager TokenManager) ValidateRefreshToken(token string) (TokenClaims, error) {
	return manager.validateToken(token, TokenTypeRefresh, "auth-refresh", []byte(manager.config.RefreshSecret))
}

func (manager TokenManager) issueToken(userID string, tokenVersion int, tokenType string, audience string, ttl time.Duration, secret []byte) (string, error) {
	now := manager.now().UTC()
	claims := TokenClaims{
		UserID:       userID,
		Type:         tokenType,
		Audience:     audience,
		Issuer:       manager.config.Issuer,
		ExpiresAt:    now.Add(ttl).Unix(),
		IssuedAt:     now.Unix(),
		NotBefore:    now.Unix(),
		TokenVersion: tokenVersion,
	}

	return signJWT(claims, secret)
}

func (manager TokenManager) validateToken(token string, tokenType string, audience string, secret []byte) (TokenClaims, error) {
	claims, err := parseJWT(token, secret)
	if err != nil {
		return TokenClaims{}, err
	}

	if claims.Type != tokenType || claims.Audience != audience || claims.Issuer != manager.config.Issuer || claims.UserID == "" {
		return TokenClaims{}, ErrInvalidToken
	}

	now := manager.now().UTC().Unix()
	if claims.NotBefore > now {
		return TokenClaims{}, ErrInvalidToken
	}
	if claims.ExpiresAt <= now {
		return TokenClaims{}, ErrExpiredToken
	}

	return claims, nil
}

func signJWT(claims TokenClaims, secret []byte) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsigned := encodedHeader + "." + encodedClaims
	signature := sign(unsigned, secret)

	return unsigned + "." + signature, nil
}

func parseJWT(token string, secret []byte) (TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return TokenClaims{}, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSignature := sign(unsigned, secret)
	if hmac.Equal([]byte(expectedSignature), []byte(parts[2])) == false {
		return TokenClaims{}, ErrInvalidToken
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return TokenClaims{}, ErrInvalidToken
	}

	var header map[string]string
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return TokenClaims{}, ErrInvalidToken
	}
	if header["alg"] != "HS256" || header["typ"] != "JWT" {
		return TokenClaims{}, ErrInvalidToken
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return TokenClaims{}, ErrInvalidToken
	}

	var claims TokenClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return TokenClaims{}, ErrInvalidToken
	}

	return claims, nil
}

func sign(unsigned string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
