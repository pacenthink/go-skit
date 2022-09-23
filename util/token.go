package util

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtValidateKey string
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
	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtValidateKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	return token, nil
}

func init() {
	jwtValidateKey = os.Getenv("JWT_VALIDATE_KEY")
}
