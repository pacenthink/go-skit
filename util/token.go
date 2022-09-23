package util

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtValidateKey string
	jwtSignKey     string
)

func ParseBearerJwtFromAuthHeader(header string) (*jwt.Token, error) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return nil, errors.New("invalid authorization header")
	}

	if strings.ToLower(parts[0]) != "bearer" {
		return nil, fmt.Errorf("unsupported token: %s", parts[0])
	}

	if parts[1] == "" {
		return nil, errors.New("empty token")
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		switch method := token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return []byte(jwtValidateKey), nil
		default:
			return nil, fmt.Errorf("unsupported signing method: %v", method)
		}
	})
}

// GenerateTokens generates and returns a HMAC token and refresh token or an error
func GenerateTokens(tokenTTL, refreshTokenTTL time.Duration) (signedToken string, signedRefreshToken string, err error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(tokenTTL)),
	}
	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSignKey))
	if err != nil {
		return
	}

	refreshClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(refreshTokenTTL)),
	}
	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtSignKey))

	return
}

func init() {
	jwtValidateKey = os.Getenv("JWT_VALIDATE_KEY")
	jwtSignKey = os.Getenv("JWT_SIGN_KEY")
}
