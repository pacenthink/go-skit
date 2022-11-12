package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var healthCheckEnv = []string{"GIT_COMMIT"}

func RegisterEnvVar(key string) {
	// We don't care about dups because these will get de-duped at runtime
	// as we use a map
	healthCheckEnv = append(healthCheckEnv, key)
}

func HealthCheckNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// HealthCheckEnvironment implements a health check that returns defined
// environment variables values registered with RegisterEnvVar if they are
// set
func HealthCheckEnvironment(c *gin.Context) {
	obj := make(map[string]string)

	for _, k := range healthCheckEnv {
		if val := os.Getenv(k); val != "" {
			obj[k] = val
		}
	}

	c.JSON(http.StatusOK, obj)
}
