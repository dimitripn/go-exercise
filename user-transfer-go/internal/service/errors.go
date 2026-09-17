package service

import "errors"

var (
	ErrUserNotFound        = errors.New("user tidak ditemukan")
	ErrTargetUserNotFound  = errors.New("target user tidak ditemukan")
	ErrInsufficientBalance = errors.New("saldo tidak mencukupi")
	ErrSelfTransfer        = errors.New("tidak bisa transfer ke diri sendiri")
	ErrInvalidNominal      = errors.New("nominal transfer harus lebih besar dari 0")
)
