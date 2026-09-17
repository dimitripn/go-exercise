package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // driver postgres untuk database/sql

	"user-transfer-go/internal/config"
)

// Connect membuka koneksi ke Postgres menggunakan konfigurasi yang diberikan
// dan memastikan koneksi benar-benar hidup lewat Ping sebelum dikembalikan.
func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi db: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingErrCh := make(chan error, 1)
	go func() { pingErrCh <- db.Ping() }()

	select {
	case err := <-pingErrCh:
		if err != nil {
			return nil, fmt.Errorf("gagal konek ke db: %w", err)
		}
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout saat konek ke db")
	}

	return db, nil
}
