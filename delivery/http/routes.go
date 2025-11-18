package http

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/config"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/repositories"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/services"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(app *fiber.App, bc services.BlockchainClient) {
	db := config.ConnectDB()

	// repos
	userRepo := repositories.NewUserRepository(db)
	claimRepo := repositories.NewClaimRepository(db)
	authRepo := repositories.NewAuthRepo(db)
	pesertaRepo := repositories.NewPesertaRepository(db)

	// services
	authService := services.NewAuthService(authRepo)
	claimSvc := services.NewClaimService(claimRepo, bc)
	pesertaService := services.NewPesertaService(pesertaRepo)

	// handlers
	authHandler := NewAuthHandler(authService)
	claimHandler := NewClaimHandler(claimSvc)
	pesertaHandler := NewPesertaHandler(pesertaService)

	api := app.Group("/api")

	api.Post("/login", authHandler.Login)
	api.Post("/register", authHandler.Register)

	claims := api.Group("/claims")
	claims.Post("/submit", claimHandler.SubmitClaim)
	claims.Post("/approve", utils.AdminOnly(), claimHandler.ApproveClaim) // protect with auth middleware in main
	claims.Get("/:id", claimHandler.GetClaim)
	claims.Get("/", claimHandler.ListClaims)

	peserta := api.Group("/peserta", utils.RequireAuth())
	peserta.Post("", pesertaHandler.CreatePeserta)
	peserta.Get("", pesertaHandler.GetAllByCreator)

	// create default admin user if none exists
	ensureAdmin(userRepo)
}

func ensureAdmin(userRepo repositories.UserRepository) {
	// create a default admin for dev if none exists
	if _, err := userRepo.FindByUsername(context.Background(), "admin"); err == nil {
		return
	}
	pass := "admin123" // dev only
	u := &models.User{
		ID:           uuid.NewString(),
		Username:     "admin",
		PasswordHash: "", // will set below
		Role:         "auditor",
	}
	h, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.PasswordHash = string(h)
	_ = userRepo.Create(context.Background(), u)
}
