package routes

import (
	"twitter-clone/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", handlers.PingHandler)
}
