package models

import "time"

// Transfer merepresentasikan satu transaksi perpindahan saldo
// dari satu user (user_id) ke user lain (target_user_id).
type Transfer struct {
	ID           int64     `json:"-"`
	UserID       int64     `json:"user_id"`
	TargetUserID int64     `json:"target_user_id"`
	Nominal      float64   `json:"nominal"`
	CreatedDate  time.Time `json:"date"`
	// Direction menjelaskan posisi user pada transfer ini ketika
	// ditampilkan lewat GET /users/{user_id}/transfers.
	// Bernilai "OUT" jika user adalah pengirim, "IN" jika user adalah penerima.
	Direction string `json:"direction,omitempty"`
}
