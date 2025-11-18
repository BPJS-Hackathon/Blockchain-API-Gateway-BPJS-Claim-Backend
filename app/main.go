package main

import (
	"log"
	"os"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/config"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/delivery/http"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/services"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// set jwt secret from env
	s := os.Getenv("JWT_SECRET")
	if s != "" {
		utils.SetSecret(s)
	}

	// connect DB
	config.ConnectDB()

	// blockchain client stub for dev
	bc := services.NewStubBlockchainClient()

	app := fiber.New()

	// register routes (this will also start background listener)
	http.RegisterRoutes(app, bc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
