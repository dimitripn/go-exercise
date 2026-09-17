package models

// CreateUserRequest adalah payload untuk membuat user baru.
type CreateUserRequest struct {
	Name    string  `json:"name" binding:"required"`
	Age     int     `json:"age" binding:"required,gt=0"`
	Balance float64 `json:"balance" binding:"gte=0"`
}

// CreateTransferRequest adalah payload untuk membuat transfer baru.
type CreateTransferRequest struct {
	TargetUserID int64   `json:"target_user_id" binding:"required"`
	Nominal      float64 `json:"nominal" binding:"required,gt=0"`
}

// ErrorResponse adalah bentuk response standar untuk error.
type ErrorResponse struct {
	Error string `json:"error"`
}
