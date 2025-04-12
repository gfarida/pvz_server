package main

import (
	"log"
	"pvz_server/internal/app/apiserver"
)

func main() {
	srv := apiserver.NewServer()

	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
