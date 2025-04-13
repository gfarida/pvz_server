package routes

import (
	"pvz_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.Engine) {
	r.POST("/dummyLogin", handlers.DummyLogin)
}
