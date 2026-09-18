package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"user-transfer-go/internal/config"
	"user-transfer-go/internal/database"
	"user-transfer-go/internal/handlers"
	"user-transfer-go/internal/repository"
	"user-transfer-go/internal/routes"
	"user-transfer-go/internal/service"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",

	Run: func(cmd *cobra.Command, args []string) {

		// --- Load config ---
		cfg := config.Load()

		// --- Database connection ---
		db, err := database.Connect(cfg)
		if err != nil {
			log.Fatalf("gagal konek ke database: %v", err)
		}
		defer db.Close()

		// --- Dependency wiring ---
		userRepo := repository.NewUserRepository(db)
		transferRepo := repository.NewTransferRepository(db)

		userService := service.NewUserService(userRepo)
		transferService := service.NewTransferService(
			userRepo,
			transferRepo,
		)

		userHandler := handlers.NewUserHandler(userService)
		transferHandler := handlers.NewTransferHandler(transferService)

		// --- Router setup ---
		router := gin.Default()

		routes.Setup(
			router,
			userHandler,
			transferHandler,
		)

		// --- Start server ---
		log.Printf("server running on :%s", cfg.AppPort)

		if err := router.Run(":" + cfg.AppPort); err != nil {
			log.Fatalf("gagal menjalankan server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
