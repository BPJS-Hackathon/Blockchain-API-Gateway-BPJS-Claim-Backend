package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB connects and auto-migrates models. Call once at startup.
func ConnectDB() *gorm.DB {
	if DB != nil {
		return DB
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// build from components if DATABASE_URL not provided
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		name := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		ssl := os.Getenv("DB_SSLMODE")
		if ssl == "" {
			ssl = "disable"
		}
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			host, user, pass, name, port, ssl)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	DB = db

	// AutoMigrate models used by Developer 3
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Claim{},
		&models.Approval{},
		&models.AuditLog{},
		&models.PesertaBPJS{},
	); err != nil {
		log.Fatal("auto migrate failed:", err)
	}

	log.Println("✅ Database connected and migrated")
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
