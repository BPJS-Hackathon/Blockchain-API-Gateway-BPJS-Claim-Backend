package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// PrintContextData mencetak semua data yang disimpan di context
func PrintContextData(c *gin.Context) {
	fmt.Printf("=== CONTEXT DATA (%s) ===\n", time.Now().Format("15:04:05"))

	// Ambil semua data dari context
	if userID, exists := c.Get("userID"); exists {
		fmt.Printf("User ID    : %s\n", userID)
	}

	if username, exists := c.Get("username"); exists {
		fmt.Printf("Username   : %s\n", username)
	}

	if role, exists := c.Get("role"); exists {
		fmt.Printf("Role       : %s\n", role)
	}

	if name, exists := c.Get("name"); exists {
		fmt.Printf("Name       : %s\n", name)
	}

	// Info request tambahan
	fmt.Printf("Method     : %s\n", c.Request.Method)
	fmt.Printf("Endpoint   : %s\n", c.Request.URL.Path)
	fmt.Printf("IP         : %s\n", c.ClientIP())
	fmt.Printf("====================================\n")
}

// PrintContextDataCompact versi compact
func PrintContextDataCompact(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	fmt.Printf("[CTX] %s %s - User: %s (%s) - Role: %s\n",
		c.Request.Method,
		c.Request.URL.Path,
		username,
		userID,
		role,
	)
}
