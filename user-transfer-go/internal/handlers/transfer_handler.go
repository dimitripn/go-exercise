package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"user-transfer-go/internal/models"
	"user-transfer-go/internal/service"
)

// TransferHandler menangani HTTP request terkait resource transfer.
type TransferHandler struct {
	transferService *service.TransferService
}

func NewTransferHandler(transferService *service.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}

// GetTransfersByUser menangani GET /users/:id/transfers
func (h *TransferHandler) GetTransfersByUser(c *gin.Context) {
	userID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "id tidak valid"})
		return
	}

	transfers, err := h.transferService.GetTransfersByUser(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "gagal mengambil riwayat transfer"})
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// CreateTransfer menangani POST /users/:id/transfers
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	userID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "id tidak valid"})
		return
	}

	var req models.CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	transfer, err := h.transferService.CreateTransfer(userID, req.TargetUserID, req.Nominal)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrTargetUserNotFound):
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrInsufficientBalance), errors.Is(err, service.ErrSelfTransfer), errors.Is(err, service.ErrInvalidNominal):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "gagal memproses transfer"})
		}
		return
	}

	c.JSON(http.StatusCreated, transfer)
}
