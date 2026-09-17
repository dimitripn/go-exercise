package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config menampung seluruh konfigurasi aplikasi yang diambil dari
// environment variable (bisa didefinisikan lewat file .env saat development).
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	AppPort    string
}

// Load membaca file .env (jika ada) lalu mengembalikan Config yang sudah
// diisi dari environment variable. Jika .env tidak ditemukan, Load tetap
// jalan (misal di production env variable di-set langsung oleh OS/container).
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "user_transfer"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
