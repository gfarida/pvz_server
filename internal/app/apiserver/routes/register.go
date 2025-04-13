package routes

import (
	"pvz_server/internal/app/deps"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, deps *deps.Dependencies) {
	registerAuthRoutes(r)
	registerPVZRoutes(r, deps)
}
