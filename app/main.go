package main

import (
	"log"
	"os"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/config"
	handlers "github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/delivery/http"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/repositories"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// Boot DB
	db, err := config.BootDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
	}

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET not found in .env")
	}

	// init repo
	authRepo := repositories.NewAuthRepo(db)
	faskes2repo := repositories.NewFaskes2Repo(db)
	adminRep := repositories.NewAdminRepo(db)

	// init services
	authService := services.NewAuthService(authRepo, jwtSecret)
	faskes2Service := services.NewFaskes2Service(faskes2repo)
	adminService := services.NewAdminService(adminRep)

	// init gin
	app := gin.Default()
	config.InitMiddleware(app)

	// init handlers
	handlers.NewAuthHandler(app, *authService, authService.GetAccessTokenManager())
	handlers.NewFaskes2Handler(app, faskes2Service, authService.GetAccessTokenManager())
	handlers.NewAdminHandler(app, adminService, authService.GetAccessTokenManager())

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	srvAddr := ":" + port

	log.Printf("🚀 Server running at http://localhost%s", srvAddr)
	if err := app.Run(srvAddr); err != nil {
		log.Fatal("❌ Failed to start server: ", err)
	}
}
