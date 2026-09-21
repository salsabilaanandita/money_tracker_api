package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// LoadEnv membaca file .env dan memuatnya sebagai environment variable
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment system")
	}
}

// ConnectDB membuat koneksi ke database PostgreSQL menggunakan GORM
func ConnectDB() {
	var dsn string

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	} else {
		sslMode := os.Getenv("DB_SSLMODE")
		if sslMode == "" {
			sslMode = "disable"
		}
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			sslMode,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Gagal konek ke database: ", err)
		return
	}

	DB = db
	log.Println("Berhasil konek ke database")
}
