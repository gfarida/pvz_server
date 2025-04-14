package routes

import (
	"pvz_server/internal/app/deps"
	"pvz_server/internal/app/middleware"
	"pvz_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func registerPVZRoutes(r *gin.Engine, deps *deps.Dependencies) {
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/pvz", handlers.CreatePVZ(deps.Store))
	protected.POST("/pvz/:pvzId/delete_last_product", handlers.DeleteLastProduct(deps.Store))
}
