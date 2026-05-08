package main

import (
	"log"
	"twitter-clone/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.RegisterRoutes(router)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
