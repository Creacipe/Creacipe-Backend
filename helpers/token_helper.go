// LOKASI: helpers/token_helper.go

package helpers

import (
	
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}


func CreateAccessToken(userID uint, roleName string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,   // Gunakan userID
		"role": roleName, // Gunakan roleName
		"exp":  time.Now().Add(time.Minute * 60).Unix(),
	})

	return claims.SignedString(getJWTSecret()) 
}


func CreateRefreshToken(userID uint, roleName string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,   // Gunakan userID
		"role": roleName, // Gunakan roleName
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return claims.SignedString(getJWTSecret()) 
}