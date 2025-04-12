package apiserver

import (
	"pvz_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/dummyLogin", handlers.DummyLogin)
}
