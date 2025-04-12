package apiserver

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()
	RegisterRoutes(r)
	return &Server{engine: r}
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}
