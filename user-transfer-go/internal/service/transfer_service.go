package service

import (
	"user-transfer-go/internal/models"
	"user-transfer-go/internal/repository"
)

// TransferService berisi business logic terkait transfer saldo antar user.
type TransferService struct {
	userRepo     *repository.UserRepository
	transferRepo *repository.TransferRepository
}

func NewTransferService(userRepo *repository.UserRepository, transferRepo *repository.TransferRepository) *TransferService {
	return &TransferService{
		userRepo:     userRepo,
		transferRepo: transferRepo,
	}
}

// GetTransfersByUser mengembalikan riwayat transfer milik user (masuk & keluar).
func (s *TransferService) GetTransfersByUser(userID int64) ([]*models.Transfer, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return s.transferRepo.FindByUserID(userID)
}

// CreateTransfer memproses transfer saldo dari userID ke targetUserID.
// Seluruh operasi (baca saldo dengan lock, validasi, update saldo, catat
// transfer) dibungkus dalam satu transaksi database supaya atomic: jika ada
// langkah yang gagal, semua perubahan di-rollback.
//
// Baris user dikunci (SELECT ... FOR UPDATE) dengan urutan id yang konsisten
// (id lebih kecil dikunci lebih dulu) untuk menghindari deadlock saat ada
// dua transfer berlawanan arah (A->B dan B->A) berjalan bersamaan.
func (s *TransferService) CreateTransfer(userID, targetUserID int64, nominal float64) (*models.Transfer, error) {
	if nominal <= 0 {
		return nil, ErrInvalidNominal
	}
	if userID == targetUserID {
		return nil, ErrSelfTransfer
	}

	tx, err := s.userRepo.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck // aman: no-op jika sudah di-Commit

	firstID, secondID := userID, targetUserID
	if secondID < firstID {
		firstID, secondID = secondID, firstID
	}

	locked := make(map[int64]*models.User, 2)
	for _, id := range []int64{firstID, secondID} {
		u, err := s.userRepo.FindByIDTx(tx, id)
		if err != nil {
			return nil, err
		}
		locked[id] = u
	}

	sender := locked[userID]
	if sender == nil {
		return nil, ErrUserNotFound
	}

	target := locked[targetUserID]
	if target == nil {
		return nil, ErrTargetUserNotFound
	}

	if sender.Balance < nominal {
		return nil, ErrInsufficientBalance
	}

	if err := s.userRepo.UpdateBalanceTx(tx, userID, -nominal); err != nil {
		return nil, err
	}
	if err := s.userRepo.UpdateBalanceTx(tx, targetUserID, nominal); err != nil {
		return nil, err
	}

	transfer, err := s.transferRepo.CreateTx(tx, userID, targetUserID, nominal)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transfer, nil
}
