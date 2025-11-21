// config/database.go
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func BootDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
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
		return nil, err
	}
	DB = db

	// AutoMigrate models
	// config/database.go - di function BootDB()
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Claims{},
		&models.RekamMedis{},
		&models.PesertaBPJS{},
		&models.DiagnosisCode{},
		&models.Faskes{}, // Jangan lupa ini
	); err != nil {
		return nil, err
	}

	// Seed data
	if err := seedAdminUser(); err != nil {
		return nil, err
	}
	if err := seedAuditorUsers(); err != nil {
		return nil, err
	}
	if err := seedFaskesUsers(); err != nil {
		return nil, err
	}
	if err := seedPesertaBPJS(); err != nil {
		return nil, err
	}
	if err := seedDiagnosisCodesFromJSON(); err != nil {
		return nil, err
	}
	log.Println("✅ Database connected and migrated with seed data")
	return DB, nil
}

func seedDiagnosisCodesFromJSON() error {
	var count int64
	DB.Model(&models.DiagnosisCode{}).Count(&count)
	if count > 0 {
		log.Println("✅ Diagnosis codes already exist, skipping seed")
		return nil
	}

	// Baca file JSON
	filePath := "data/diagnosis_seed.json"
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read diagnosis codes JSON file: %v", err)
	}

	var diagnosisData []models.DiagnosisCode
	if err := json.Unmarshal(fileData, &diagnosisData); err != nil {
		return fmt.Errorf("failed to unmarshal diagnosis codes JSON: %v", err)
	}

	// Create semua diagnosis codes
	for i := range diagnosisData {
		if err := DB.Create(&diagnosisData[i]).Error; err != nil {
			return fmt.Errorf("failed to create diagnosis code %s: %v", diagnosisData[i].Code, err)
		}
	}

	log.Printf("✅ Seeded %d diagnosis codes from JSON", len(diagnosisData))
	return nil
}

// Seed PesertaBPJS data dengan looping
// config/database.go - update seedPesertaBPJS function
func seedPesertaBPJS() error {
	var count int64
	DB.Model(&models.PesertaBPJS{}).Count(&count)
	if count > 0 {
		log.Println("✅ PesertaBPJS data already exists, skipping seed")
		return nil
	}

	// Dapatkan ID faskes dari database
	var faskes1 models.Faskes
	var faskes2 models.Faskes

	if err := DB.Where("nama_faskes = ?", getEnvOrDefault("FASKES1NAME", "Faskes Pertama")).First(&faskes1).Error; err != nil {
		return fmt.Errorf("failed to find faskes1: %v", err)
	}
	if err := DB.Where("nama_faskes = ?", getEnvOrDefault("FASKES2NAME", "Faskes Kedua")).First(&faskes2).Error; err != nil {
		return fmt.Errorf("failed to find faskes2: %v", err)
	}

	// Data peserta untuk Faskes 1 (Puskesmas)
	pesertaFaskes1 := []models.PesertaBPJS{
		{
			PSTV01: "0001234567890",
			PSTV02: "FAM001",
			PSTV03: time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC),
			PSTV04: 1,
			PSTV05: 1,
			PSTV06: 2,
			PSTV07: "III",
			PSTV08: 2,
			PSTV09: "31",
			PSTV10: "3171",
			PSTV11: faskes1.NamaFaskes, // "Faskes Pertama"
			PSTV12: 1,                  // Puskesmas
			PSTV13: "31",
			PSTV14: "3171",
			PSTV15: 1.0000,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
		{
			PSTV01: "0001234567891",
			PSTV02: "FAM001",
			PSTV03: time.Date(1987, 8, 20, 0, 0, 0, 0, time.UTC),
			PSTV04: 2,
			PSTV05: 2,
			PSTV06: 2,
			PSTV07: "III",
			PSTV08: 2,
			PSTV09: "31",
			PSTV10: "3171",
			PSTV11: faskes1.NamaFaskes, // "Faskes Pertama"
			PSTV12: 1,                  // Puskesmas
			PSTV13: "31",
			PSTV14: "3171",
			PSTV15: 1.0000,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
		{
			PSTV01: "0001234567892",
			PSTV02: "FAM001",
			PSTV03: time.Date(2015, 11, 10, 0, 0, 0, 0, time.UTC),
			PSTV04: 4,
			PSTV05: 1,
			PSTV06: 1,
			PSTV07: "III",
			PSTV08: 2,
			PSTV09: "31",
			PSTV10: "3171",
			PSTV11: faskes1.NamaFaskes, // "Faskes Pertama"
			PSTV12: 1,                  // Puskesmas
			PSTV13: "31",
			PSTV14: "3171",
			PSTV15: 0.7500,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
		{
			PSTV01: "0001234567896",
			PSTV02: "FAM004",
			PSTV03: time.Date(1965, 1, 25, 0, 0, 0, 0, time.UTC),
			PSTV04: 1,
			PSTV05: 2,
			PSTV06: 3,
			PSTV07: "III",
			PSTV08: 3,
			PSTV09: "34",
			PSTV10: "3471",
			PSTV11: faskes1.NamaFaskes, // "Faskes Pertama"
			PSTV12: 1,                  // Puskesmas
			PSTV13: "34",
			PSTV14: "3471",
			PSTV15: 1.0000,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
	}

	// Data peserta untuk Faskes 2 (Klinik)
	pesertaFaskes2 := []models.PesertaBPJS{
		{
			PSTV01: "0001234567893",
			PSTV02: "FAM002",
			PSTV03: time.Date(1978, 12, 5, 0, 0, 0, 0, time.UTC),
			PSTV04: 1,
			PSTV05: 1,
			PSTV06: 2,
			PSTV07: "II",
			PSTV08: 1,
			PSTV09: "32",
			PSTV10: "3273",
			PSTV11: faskes2.NamaFaskes, // "Faskes Kedua"
			PSTV12: 2,                  // Klinik
			PSTV13: "32",
			PSTV14: "3273",
			PSTV15: 1.2500,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
		{
			PSTV01: "0001234567894",
			PSTV02: "FAM002",
			PSTV03: time.Date(1980, 4, 18, 0, 0, 0, 0, time.UTC),
			PSTV04: 2,
			PSTV05: 2,
			PSTV06: 2,
			PSTV07: "II",
			PSTV08: 1,
			PSTV09: "32",
			PSTV10: "3273",
			PSTV11: faskes2.NamaFaskes, // "Faskes Kedua"
			PSTV12: 2,                  // Klinik
			PSTV13: "32",
			PSTV14: "3273",
			PSTV15: 1.2500,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
		{
			PSTV01: "0001234567897",
			PSTV02: "FAM005",
			PSTV03: time.Date(1972, 9, 14, 0, 0, 0, 0, time.UTC),
			PSTV04: 1,
			PSTV05: 1,
			PSTV06: 2,
			PSTV07: "II",
			PSTV08: 1,
			PSTV09: "35",
			PSTV10: "3578",
			PSTV11: faskes2.NamaFaskes, // "Faskes Kedua"
			PSTV12: 2,                  // Klinik
			PSTV13: "35",
			PSTV14: "3578",
			PSTV15: 1.2500,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
		{
			PSTV01: "0001234567898",
			PSTV02: "FAM005",
			PSTV03: time.Date(1975, 6, 8, 0, 0, 0, 0, time.UTC),
			PSTV04: 2,
			PSTV05: 2,
			PSTV06: 2,
			PSTV07: "II",
			PSTV08: 1,
			PSTV09: "35",
			PSTV10: "3578",
			PSTV11: faskes2.NamaFaskes, // "Faskes Kedua"
			PSTV12: 2,                  // Klinik
			PSTV13: "35",
			PSTV14: "3578",
			PSTV15: 1.2500,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
		{
			PSTV01: "0001234567899",
			PSTV02: "FAM006",
			PSTV03: time.Date(1988, 2, 28, 0, 0, 0, 0, time.UTC),
			PSTV04: 1,
			PSTV05: 1,
			PSTV06: 1,
			PSTV07: "I",
			PSTV08: 1,
			PSTV09: "36",
			PSTV10: "3671",
			PSTV11: faskes2.NamaFaskes, // "Faskes Kedua"
			PSTV12: 2,                  // Klinik
			PSTV13: "36",
			PSTV14: "3671",
			PSTV15: 1.5000,
			PSTV16: 2024,
			PSTV17: "Aktif",
			PSTV18: nil,
		},
	}

	// Gabungkan semua peserta
	allPeserta := append(pesertaFaskes1, pesertaFaskes2...)

	// Create semua peserta dalam batch
	if err := DB.CreateInBatches(allPeserta, 10).Error; err != nil {
		return fmt.Errorf("failed to create peserta BPJS: %v", err)
	}

	log.Printf("✅ Seeded %d PesertaBPJS records", len(allPeserta))
	log.Printf("   - %d peserta for %s", len(pesertaFaskes1), faskes1.NamaFaskes)
	log.Printf("   - %d peserta for %s", len(pesertaFaskes2), faskes2.NamaFaskes)
	return nil
}

// Fungsi seed lainnya tetap sama...
func seedAdminUser() error {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		adminPass := os.Getenv("ADMIN_PASSWORD")
		if adminPass == "" {
			adminPass = "admin123"
		}
		adminHashed, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)

		adminUser := models.User{
			Name:         "Admin",
			Username:     getEnvOrDefault("ADMIN_USERNAME", "admin"),
			PasswordHash: string(adminHashed),
			Role:         models.RoleAdmin,
		}
		DB.Create(&adminUser)
		log.Println("✅ Admin user seeded")
	}
	return nil
}

func seedAuditorUsers() error {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleAuditor).Count(&count)
	if count == 0 {
		adminPass := os.Getenv("AUDITOR_PASSWORD")
		if adminPass == "" {
			adminPass = "auditor123"
		}
		adminHashed, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)

		adminUser := models.User{
			Name:         "Tom Brady The Auditor",
			Username:     getEnvOrDefault("AUDITOR_USERNAME", "auditor"),
			PasswordHash: string(adminHashed),
			Role:         models.RoleAuditor,
		}
		DB.Create(&adminUser)
		log.Println("✅ Auditor user seeded")
	}
	return nil
}

// config/database.go
// config/database.go
func seedFaskesUsers() error {
	// Check if faskes users already exist
	var userCount int64
	DB.Model(&models.User{}).Where("role = ?", "faskes").Count(&userCount)
	if userCount > 0 {
		log.Println("✅ Faskes users already exist, skipping seed")
		return nil
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		// Seed Faskes data - BIARKAN ID AUTO GENERATE
		faskes1 := models.Faskes{
			NamaFaskes:  getEnvOrDefault("FASKES1NAME", "Faskes Pertama"),
			JenisFaskes: getEnvOrDefault("JENIS_FASKES1", "FASKES 1"),
		}
		faskes2 := models.Faskes{
			NamaFaskes:  getEnvOrDefault("FASKES2NAME", "Faskes Kedua"),
			JenisFaskes: getEnvOrDefault("JENIS_FASKES2", "FASKES 2"),
		}

		// Create faskes dan ambil ID yang di-generate
		if err := tx.Create(&faskes1).Error; err != nil {
			return fmt.Errorf("failed to create faskes1: %v", err)
		}
		if err := tx.Create(&faskes2).Error; err != nil {
			return fmt.Errorf("failed to create faskes2: %v", err)
		}

		// Seed Users dengan reference ke Faskes
		faskes1Hashed, _ := bcrypt.GenerateFromPassword([]byte(getEnvOrDefault("FASKES1PASS", "faskes123")), bcrypt.DefaultCost)
		faskes2Hashed, _ := bcrypt.GenerateFromPassword([]byte(getEnvOrDefault("FASKES2PASS", "faskes123")), bcrypt.DefaultCost)

		user1 := models.User{
			Name:         "Tom",
			Username:     getEnvOrDefault("FASKES1", "faskes1"),
			PasswordHash: string(faskes1Hashed),
			Role:         models.RoleFaskes,
			FaskesID:     &faskes1.ID, // Gunakan ID yang sudah di-generate
		}
		user2 := models.User{
			Name:         "Brady",
			Username:     getEnvOrDefault("FASKES2", "faskes2"),
			PasswordHash: string(faskes2Hashed),
			Role:         models.RoleFaskes,
			FaskesID:     &faskes2.ID, // Gunakan ID yang sudah di-generate
		}

		if err := tx.Create(&user1).Error; err != nil {
			return fmt.Errorf("failed to create faskes1 user: %v", err)
		}
		if err := tx.Create(&user2).Error; err != nil {
			return fmt.Errorf("failed to create faskes2 user: %v", err)
		}

		log.Printf("✅ Seeded 2 faskes and 2 faskes users")
		log.Printf("   - %s (ID: %s)", faskes1.NamaFaskes, faskes1.ID)
		log.Printf("   - %s (ID: %s)", faskes2.NamaFaskes, faskes2.ID)
		return nil
	})
}

// Helper function untuk environment variables
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
