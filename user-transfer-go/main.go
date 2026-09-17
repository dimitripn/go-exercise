package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"user-transfer-go/internal/config"
	"user-transfer-go/internal/database"
	"user-transfer-go/internal/handlers"
	"user-transfer-go/internal/repository"
	"user-transfer-go/internal/routes"
	"user-transfer-go/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("gagal konek ke database: %v", err)
	}
	defer db.Close()

	// --- Dependency wiring (repository -> service -> handler) ---
	userRepo := repository.NewUserRepository(db)
	transferRepo := repository.NewTransferRepository(db)

	userService := service.NewUserService(userRepo)
	transferService := service.NewTransferService(userRepo, transferRepo)

	userHandler := handlers.NewUserHandler(userService)
	transferHandler := handlers.NewTransferHandler(transferService)

	// --- Router setup ---
	router := gin.Default()
	routes.Setup(router, userHandler, transferHandler)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
