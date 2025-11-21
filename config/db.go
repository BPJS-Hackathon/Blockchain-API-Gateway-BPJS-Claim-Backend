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
func seedPesertaBPJS() error {
	var count int64
	DB.Model(&models.PesertaBPJS{}).Count(&count)
	if count > 0 {
		log.Println("✅ PesertaBPJS data already exists, skipping seed")
		return nil
	}

	// Data base untuk generate peserta
	families := []struct {
		familyID string
		members  []struct {
			gender     int
			birthYear  int
			birthMonth time.Month
			birthDay   int
			relation   int
			marital    int
		}
		province   string
		city       string
		faskes     string
		faskesType int
		class      string
		segment    int
		bobot      float64
	}{
		// Keluarga 1 - Jakarta
		{
			familyID: "FAM001",
			members: []struct {
				gender     int
				birthYear  int
				birthMonth time.Month
				birthDay   int
				relation   int
				marital    int
			}{
				{1, 1985, 3, 15, 1, 2},  // Kepala keluarga (laki-laki, menikah)
				{2, 1987, 8, 20, 2, 2},  // Istri
				{1, 2015, 11, 10, 4, 1}, // Anak laki-laki
			},
			province:   "31",
			city:       "3171",
			faskes:     "Puskesmas Jakarta Selatan",
			faskesType: 1,
			class:      "III",
			segment:    2,
			bobot:      1.0000,
		},
		// Keluarga 2 - Bandung
		{
			familyID: "FAM002",
			members: []struct {
				gender     int
				birthYear  int
				birthMonth time.Month
				birthDay   int
				relation   int
				marital    int
			}{
				{1, 1978, 12, 5, 1, 2},
				{2, 1980, 4, 18, 2, 2},
			},
			province:   "32",
			city:       "3273",
			faskes:     "Klinik Bandung",
			faskesType: 2,
			class:      "II",
			segment:    1,
			bobot:      1.2500,
		},
		// Keluarga 3 - Solo (single)
		{
			familyID: "FAM003",
			members: []struct {
				gender     int
				birthYear  int
				birthMonth time.Month
				birthDay   int
				relation   int
				marital    int
			}{
				{1, 1992, 7, 30, 1, 1},
			},
			province:   "33",
			city:       "3372",
			faskes:     "RSUD Solo",
			faskesType: 3,
			class:      "I",
			segment:    1,
			bobot:      1.5000,
		},
		// Keluarga 4 - Yogyakarta (single parent)
		{
			familyID: "FAM004",
			members: []struct {
				gender     int
				birthYear  int
				birthMonth time.Month
				birthDay   int
				relation   int
				marital    int
			}{
				{2, 1965, 1, 25, 1, 3},
			},
			province:   "34",
			city:       "3471",
			faskes:     "Puskesmas Yogyakarta",
			faskesType: 1,
			class:      "III",
			segment:    3,
			bobot:      1.0000,
		},
		// Keluarga 5 - Surabaya
		{
			familyID: "FAM005",
			members: []struct {
				gender     int
				birthYear  int
				birthMonth time.Month
				birthDay   int
				relation   int
				marital    int
			}{
				{1, 1972, 9, 14, 1, 2},
				{2, 1975, 6, 8, 2, 2},
				{2, 2010, 2, 28, 4, 1}, // Anak perempuan
			},
			province:   "35",
			city:       "3578",
			faskes:     "RS Surabaya",
			faskesType: 3,
			class:      "II",
			segment:    1,
			bobot:      1.2500,
		},
	}

	// Generate peserta berdasarkan families
	pesertaNumber := 890
	var allPeserta []models.PesertaBPJS

	for _, family := range families {
		for _, member := range family.members {
			pesertaNumber++

			// Adjust bobot untuk anak
			bobot := family.bobot
			if member.relation == 4 { // Anak
				bobot = 0.7500
			}

			peserta := models.PesertaBPJS{
				PSTV01: fmt.Sprintf("0001234567%d", pesertaNumber),
				PSTV02: family.familyID,
				PSTV03: time.Date(member.birthYear, member.birthMonth, member.birthDay, 0, 0, 0, 0, time.UTC),
				PSTV04: member.relation,
				PSTV05: member.gender,
				PSTV06: member.marital,
				PSTV07: family.class,
				PSTV08: family.segment,
				PSTV09: family.province,
				PSTV10: family.city,
				PSTV11: family.faskes,
				PSTV12: family.faskesType,
				PSTV13: family.province,
				PSTV14: family.city,
				PSTV15: bobot,
				PSTV16: 2024,
				PSTV17: "Aktif",
				PSTV18: nil,
			}
			allPeserta = append(allPeserta, peserta)
		}
	}

	// Create semua peserta dalam batch
	if err := DB.CreateInBatches(allPeserta, 10).Error; err != nil {
		return fmt.Errorf("failed to create peserta BPJS: %v", err)
	}

	log.Printf("✅ Seeded %d PesertaBPJS records", len(allPeserta))
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
			Username:     getEnvOrDefault("ADMIN_USERNAME", "admin"),
			PasswordHash: string(adminHashed),
			Role:         models.RoleAdmin,
		}
		DB.Create(&adminUser)
		log.Println("✅ Admin user seeded")
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
			Username:     getEnvOrDefault("FASKES1", "faskes1"),
			PasswordHash: string(faskes1Hashed),
			Role:         models.RoleFaskes,
			FaskesID:     &faskes1.ID, // Gunakan ID yang sudah di-generate
		}
		user2 := models.User{
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
