package repository

import (
	"database/sql"

	"user-transfer-go/internal/models"
)

// TransferRepository menangani akses data transfer ke tabel "transfers" di Postgres.
type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

// CreateTx mencatat transfer baru di dalam sebuah transaksi *sql.Tx.
func (r *TransferRepository) CreateTx(tx *sql.Tx, userID, targetUserID int64, nominal float64) (*models.Transfer, error) {
	t := &models.Transfer{
		UserID:       userID,
		TargetUserID: targetUserID,
		Nominal:      nominal,
	}
	err := tx.QueryRow(
		`INSERT INTO transfers (user_id, target_user_id, nominal)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_date`,
		userID, targetUserID, nominal,
	).Scan(&t.ID, &t.CreatedDate)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// FindByUserID mengembalikan seluruh transfer yang melibatkan user tersebut,
// baik sebagai pengirim maupun penerima, diurutkan terbaru lebih dulu.
// Field Direction diisi "OUT" jika user adalah pengirim, "IN" jika penerima.
func (r *TransferRepository) FindByUserID(userID int64) ([]*models.Transfer, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, target_user_id, nominal, created_date,
		       CASE WHEN user_id = $1 THEN 'OUT' ELSE 'IN' END AS direction
		FROM transfers
		WHERE user_id = $1 OR target_user_id = $1
		ORDER BY created_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := make([]*models.Transfer, 0)
	for rows.Next() {
		t := &models.Transfer{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.TargetUserID, &t.Nominal, &t.CreatedDate, &t.Direction); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)
	}
	return transfers, rows.Err()
}
