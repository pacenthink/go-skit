package util

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	defaultTokenTTL        = 24 * time.Hour
	defaultRefreshTokenTTL = 168 * time.Hour
)

type TokenPair struct {
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type Claims struct {
	// Identity provider e.g. github, gitlab, bitbucket etc.
	IDP string `json:"idp"`

	// Username based on idp
	Alias string `json:"alias"`

	// Standard JWT claims
	*jwt.RegisteredClaims
}

func NewClaims() *Claims {
	return &Claims{
		RegisteredClaims: &jwt.RegisteredClaims{
			// Issuer: ,
			// IssuedAt: ,
			// ID: "Account.ID",
			// NotBefore: ,
			// ExpiresAt: ,
			// Audience: ,
			// Subject: ,
		},
	}
}

// NewPairWithClaims creates a new token and refresh token pair
func NewTokenPairWithClaims(c *Claims, signMethod jwt.SigningMethod) (*TokenPair, error) {
	jwtSignKey := os.Getenv("JWT_SIGN_KEY")
	if jwtSignKey == "" {
		return nil, errors.New("sign key not set")
	}

	c.IssuedAt = jwt.NewNumericDate(time.Now().Local())
	c.ExpiresAt = jwt.NewNumericDate(time.Now().Local().Add(defaultTokenTTL))
	token, err := jwt.NewWithClaims(signMethod, c).SignedString([]byte(jwtSignKey))
	if err != nil {
		return nil, err
	}

	c.IssuedAt = jwt.NewNumericDate(time.Now().Local())
	c.ExpiresAt = jwt.NewNumericDate(time.Now().Local().Add(defaultRefreshTokenTTL))
	refresh, err := jwt.NewWithClaims(signMethod, c).SignedString([]byte(jwtSignKey))
	if err != nil {
		return nil, err
	}

	return &TokenPair{token, refresh}, nil
}
