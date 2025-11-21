package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitMiddleware(app *gin.Engine) {
	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("ALLOW_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Logging & Recovery
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	// Security Headers
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Next()
	})

	// Timeout Middleware
	app.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-ctx.Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
			c.Abort()
		case <-done:
			// selesai normal
		}
	})
}

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Need Credential to Access this Resource (Authorization header missing)",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Ambil token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userUUID, name, username, role, err := jwtManager.VerifyToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Simpan ke context
		c.Set("userID", userUUID)
		c.Set("username", username) // Sekarang username tersedia
		c.Set("role", role)
		c.Set("name", name)
		// Lanjut ke handler berikutnya
		c.Next()
	}
}

func Faskes2Only() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Username not found in token",
			})
			c.Abort()
			return
		}

		// Check semua kemungkinan username untuk faskes2
		possibleFaskes2Names := []string{"faskes2", "FASKES2", "faskes 2", "Faskes2"}
		usernameStr := username.(string)

		for _, name := range possibleFaskes2Names {
			if usernameStr == name {
				fmt.Printf("Match found with: %s\n", name)
				c.Next()
				return
			}
		}

		fmt.Printf("No match found. User '%s' is not Faskes2\n", usernameStr)
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Access restricted to Faskes2 users only. Current user: " + usernameStr,
		})
		c.Abort()
	}
}

func Faskes1Only() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Username not found in token",
			})
			c.Abort()
			return
		}

		// Check semua kemungkinan username untuk faskes2
		possibleFaskes2Names := []string{"faskes1", "FASKES1", "faskes 1", "Faskes1"}
		usernameStr := username.(string)

		for _, name := range possibleFaskes2Names {
			if usernameStr == name {
				fmt.Printf("Match found with: %s\n", name)
				c.Next()
				return
			}
		}

		fmt.Printf("No match found. User '%s' is not Faskes1\n", usernameStr)
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Access restricted to Faskes1 users only. Current user: " + usernameStr,
		})
		c.Abort()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Username not found in token",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access restricted to Admin users only",
			})
		}

		c.Next()
	}
}
