package routes

import (
	"github.com/gin-gonic/gin"

	"user-transfer-go/internal/handlers"
)

// Setup mendaftarkan seluruh route API ke router gin.
func Setup(router *gin.Engine, userHandler *handlers.UserHandler, transferHandler *handlers.TransferHandler) {
	users := router.Group("/users")
	{
		users.GET("", userHandler.GetUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUserByID)
		users.DELETE("/:id", userHandler.DeleteUser)

		users.GET("/:id/transfers", transferHandler.GetTransfersByUser)
		users.POST("/:id/transfers", transferHandler.CreateTransfer)
	}
}
