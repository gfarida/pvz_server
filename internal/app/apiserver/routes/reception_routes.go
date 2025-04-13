package routes

import (
	"pvz_server/internal/app/deps"
	"pvz_server/internal/app/middleware"
	"pvz_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func registerReceptionRoutes(r *gin.Engine, d *deps.Dependencies) {
	protected := r.Group(("/"))
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/receprions", handlers.CreateReceprion(d.Store))
}
