// LOKASI: helpers/token_helper.go

package helpers

import (
	// "creacipe-backend/models" // <-- Sudah tidak perlu import models
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PERBAIKAN: Jangan inisialisasi jwtSecret di sini!
// var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // HAPUS INI

/**
 * getJWTSecret mengambil JWT secret dari environment variable
 * Dipanggil setiap kali dibutuhkan untuk memastikan nilai terbaru
 */
func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

/**
 * CreateAccessToken membuat token akses jangka pendek.
 * Terima UserID dan RoleName secara EKSPLISIT
 */
func CreateAccessToken(userID uint, roleName string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,   // Gunakan userID
		"role": roleName, // Gunakan roleName
		"exp":  time.Now().Add(time.Minute * 60).Unix(),
	})

	return claims.SignedString(getJWTSecret()) // Gunakan fungsi getJWTSecret()
}

/**
 * CreateRefreshToken membuat token refresh jangka panjang.
 * Terima UserID dan RoleName secara EKSPLISIT
 */
func CreateRefreshToken(userID uint, roleName string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,   // Gunakan userID
		"role": roleName, // Gunakan roleName
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return claims.SignedString(getJWTSecret()) // Gunakan fungsi getJWTSecret()
}