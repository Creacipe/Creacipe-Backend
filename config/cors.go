// config/cors.go
package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	// Ganti dengan alamat frontend Anda saat production
	config.AllowOrigins = []string{"http://localhost:5173"} 
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	return cors.New(config)
}