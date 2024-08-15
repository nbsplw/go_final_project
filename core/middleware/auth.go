package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"main/core/config"
	"time"
)

var secret []byte

func Init() {
	secret = []byte(config.Get().Auth.Key)
}

type Claims struct {
	Hash string `json:"hash"`
	jwt.Claims
}

func AuthMiddleware(c *fiber.Ctx) error {
	if config.Get().Auth.Password != "" {
		token := c.Cookies("token")
		if token == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		if ok, _ := ValidateToken(token); !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	return c.Next()
}

func GenerateToken(hash string) (string, error) {
	exp := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		Hash: hash,
		Claims: jwt.MapClaims{
			"exp": exp,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateToken(t string) (bool, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}
