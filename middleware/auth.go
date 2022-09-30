package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pacenthink/go-skit/token"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		tkn, err := token.ParseBearerJwtFromAuthHeader(authHeader)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !tkn.Valid {
			c.AbortWithError(http.StatusUnauthorized, errors.New("token invalid"))
			return
		}

		if err = tkn.Claims.Valid(); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("token", tkn)
		c.Next()
	}
}
