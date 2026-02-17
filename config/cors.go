// config/cors.go
package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://127.0.0.1:5500",
		"http://localhost:5500",
		"https://creacipe.vercel.app",
		
	}

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	config.AllowCredentials = true


	return cors.New(config)
}