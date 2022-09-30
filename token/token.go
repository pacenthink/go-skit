package token

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
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
	return ValidateToken(parts[1])
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		switch method := token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			jwtValidateKey := os.Getenv("JWT_VALIDATE_KEY")
			if jwtValidateKey == "" {
				return nil, errors.New("validate key not set")
			}
			return []byte(jwtValidateKey), nil

		default:
			return nil, fmt.Errorf("unsupported signing method: %v", method)

		}
	})
}
