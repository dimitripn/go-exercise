package repository

import (
	"database/sql"
	"errors"

	"user-transfer-go/internal/models"
)

// UserRepository menangani akses data user ke tabel "users" di Postgres.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll mengembalikan seluruh user, diurutkan berdasarkan id.
func (r *UserRepository) FindAll() ([]*models.User, error) {
	rows, err := r.db.Query(`SELECT id, name, age, balance FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Balance); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// FindByID mencari user berdasarkan id. Mengembalikan (nil, nil) jika tidak ditemukan.
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, name, age, balance FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Age, &u.Balance)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Create menyimpan user baru. Id di-generate otomatis oleh Postgres (SERIAL).
func (r *UserRepository) Create(name string, age int, balance float64) (*models.User, error) {
	u := &models.User{Name: name, Age: age, Balance: balance}
	err := r.db.QueryRow(
		`INSERT INTO users (name, age, balance) VALUES ($1, $2, $3) RETURNING id`,
		name, age, balance,
	).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Delete menghapus user berdasarkan id. Mengembalikan false jika user tidak ditemukan.
func (r *UserRepository) Delete(id int64) (bool, error) {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

// FindByIDTx sama seperti FindByID tapi dijalankan di dalam transaksi dan
// mengunci baris (FOR UPDATE), dipakai sebelum validasi & update saldo saat
// transfer agar aman dari race condition antar transfer yang berjalan bersamaan.
func (r *UserRepository) FindByIDTx(tx *sql.Tx, id int64) (*models.User, error) {
	u := &models.User{}
	err := tx.QueryRow(
		`SELECT id, name, age, balance FROM users WHERE id = $1 FOR UPDATE`, id,
	).Scan(&u.ID, &u.Name, &u.Age, &u.Balance)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateBalanceTx mengubah saldo user (delta bisa positif/negatif) di dalam
// sebuah transaksi *sql.Tx, dipakai saat proses transfer agar atomic.
func (r *UserRepository) UpdateBalanceTx(tx *sql.Tx, id int64, delta float64) error {
	_, err := tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, delta, id)
	return err
}

// BeginTx memulai transaksi baru. Dipakai oleh service layer untuk
// membungkus proses transfer (baca saldo, validasi, update, catat) secara atomic.
func (r *UserRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
