package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"user-transfer-go/internal/models"
	"user-transfer-go/internal/service"
)

// UserHandler menangani HTTP request terkait resource user.
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers menangani GET /users
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "gagal mengambil data user"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID menangani GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "id tidak valid"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "gagal mengambil data user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser menangani POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Age, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "gagal membuat user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// DeleteUser menangani DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "id tidak valid"})
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "gagal menghapus user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}
