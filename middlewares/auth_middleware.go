// middlewares/auth_middleware.go
package middlewares

import (
	"creacipe-backend/config" // Ganti dengan nama modul Anda
	"creacipe-backend/models" // Ganti dengan nama modul Anda
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth adalah middleware untuk memeriksa token JWT yang valid.
func RequireAuth(c *gin.Context) {
	// 1. Ambil header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header otorisasi tidak ditemukan"})
		return
	}

	// 2. Periksa format "Bearer [token]"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token harus 'Bearer [token]'"})
		return
	}

	// 3. Parse dan validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Metode signing tidak terduga: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// 4. Ambil claims dan periksa apakah token valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 5. Periksa apakah token sudah kedaluwarsa
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token sudah kedaluwarsa"})
			return
		}

		// 6. Ambil data user dari database
		var user models.User
		userID := uint(claims["sub"].(float64))
		config.DB.First(&user, userID)

		if user.UserID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User pemilik token tidak ditemukan"})
			return
		}

		// 7. Simpan data user di context Gin untuk digunakan di controller
		c.Set("user", user)

		// Lanjutkan ke handler berikutnya
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
	}
}

// AuthorizeRole adalah middleware untuk memeriksa peran pengguna.
func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil data user yang sudah disimpan oleh RequireAuth
		userInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan di context"})
			return
		}

		// Konversi interface ke struct User
		user := userInterface.(models.User)
		
		isAllowed := false
		// Periksa apakah peran user ada di dalam daftar peran yang diizinkan
		for _, role := range allowedRoles {
			if user.Role == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: peran tidak diizinkan"})
			return
		}

		// Lanjutkan jika peran diizinkan
		c.Next()
	}
}