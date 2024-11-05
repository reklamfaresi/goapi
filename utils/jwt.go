package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("secret_key") // Burada gizli bir anahtar tanımlıyoruz (daha güvenli bir anahtar kullanın)

// JWT oluşturma fonksiyonu
func GenerateJWT(username, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token süresi: 24 saat
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Kullanıcı rolünü almak için JWT doğrulama fonksiyonu
func GetUserRoleFromJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role := claims["role"].(string)
		return role, nil
	}

	return "", fmt.Errorf("Geçersiz token")
}

// JWT'den kullanıcı adını almak için fonksiyon
func GetUsernameFromJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}

	return "", fmt.Errorf("Geçersiz token")
}
