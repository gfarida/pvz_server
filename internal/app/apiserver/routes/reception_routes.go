package routes

import (
	"pvz_server/internal/app/deps"
	"pvz_server/internal/app/middleware"
	"pvz_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func registerReceptionRoutes(r *gin.Engine, deps *deps.Dependencies) {
	protected := r.Group(("/"))
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/receptions", handlers.CreateReception(deps.Store))
	protected.POST("/pvz/:pvzId/close_last_reception", handlers.CloseLastReception(deps.Store))
}
