package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pacenthink/go-skit/util"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		token, err := util.ParseBearerJwtFromAuthHeader(authHeader)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, errors.New("token invalid"))
			return
		}

		if err = token.Claims.Valid(); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("token", token)
		c.Next()
	}
}
