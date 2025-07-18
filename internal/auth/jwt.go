package auth

import (
	"crypto/rand"
	"github.com/Laky-64/gologging"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

var jwtSecret []byte

func init() {
	jwtSecret = make([]byte, 32)
	_, err := rand.Read(jwtSecret)
	if err != nil {
		gologging.FatalF("Error generating random secret: %s", err)
	}
}

func GenerateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(int(userID)),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
