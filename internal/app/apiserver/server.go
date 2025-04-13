package apiserver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"pvz_server/internal/app/apiserver/routes"
	"pvz_server/internal/app/deps"
	"pvz_server/internal/app/store"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	db, err := connectDB()

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	s := &Server{
		engine: gin.Default(),
	}

	deps := &deps.Dependencies{
		Store: store.New(db),
	}

	routes.RegisterRoutes(s.engine, deps)
	return s
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func connectDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	return sql.Open("postgres", dsn)
}
