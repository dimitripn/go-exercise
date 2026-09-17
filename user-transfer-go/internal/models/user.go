package models

// User merepresentasikan data pengguna beserta saldo miliknya.
type User struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Balance float64 `json:"balance"`
}
