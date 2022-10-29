package middleware

import (
	"praktikum/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userID int, name string) (string, error) {
	claims := jwt.MapClaims{}
	claims["userId"] = userID
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Encoder
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//return token + error message
	return token.SignedString([]byte(config.TokenSecret))
}
