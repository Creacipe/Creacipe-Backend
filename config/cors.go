// config/cors.go
package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	// Memulai dengan konfigurasi default yang aman
	config := cors.DefaultConfig()

	// --- SESUAIKAN DI SINI ---
	// Tambahkan semua alamat frontend yang ingin Anda izinkan.
	// Ini bisa alamat lokal untuk pengembangan atau alamat domain saat produksi.
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://127.0.0.1:5500",
		"http://localhost:5500",
		"https://creacipe.vercel.app",
		
	}

	// Pastikan semua metode HTTP yang Anda butuhkan diizinkan
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

	// Izinkan header 'Authorization' agar frontend bisa mengirim token JWT
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	// Konfigurasi lain (opsional)
	config.AllowCredentials = true

	// Terapkan konfigurasi
	return cors.New(config)
}