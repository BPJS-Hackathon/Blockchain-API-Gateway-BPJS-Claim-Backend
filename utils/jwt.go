package utils

import (
	"fmt"
	"time"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func GenerateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
}

// Middleware to require auth and set user_id and role in Locals
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing auth"})
		}
		// expect "Bearer <token>"
		var tokenStr string
		if _, err := fmt.Sscanf(auth, "Bearer %s", &tokenStr); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid auth header"})
		}
		tok, err := ParseToken(tokenStr)
		if err != nil || !tok.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		if claims, ok := tok.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				c.Locals("user_id", sub)
			}
			if role, ok := claims["role"].(string); ok {
				c.Locals("role", role)
			}
		}
		return c.Next()
	}
}

// Middleware admin role require
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != models.RoleAdmin {
			return c.Status(403).JSON(fiber.Map{"error": "access denied"})
		}

		return c.Next()
	}
}

func SetSecret(secret string) {
	jwtSecret = []byte(secret)
}
