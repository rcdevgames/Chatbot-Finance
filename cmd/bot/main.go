package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"chatbot/internal/config"
	"chatbot/internal/database"
	"chatbot/internal/handler"
	"chatbot/internal/repository"
	"chatbot/internal/service"
	"chatbot/pkg/telegram"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize dependencies
	telegramClient := telegram.NewClient(cfg.TelegramBotToken)
	groqService := service.NewGroqService(cfg.GroqAPIKey)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)
	chatHistoryRepo := repository.NewChatHistoryRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)

	// Initialize services
	userService := service.NewUserService(userRepo, transactionRepo, chatHistoryRepo, categoryRepo)
	licenseService := service.NewLicenseService(db.DB)

	// Initialize handlers
	telegramHandler := handler.NewTelegramHandler(telegramClient, userService, groqService, licenseService)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		health := gin.H{
			"status": "ok",
			"message": "Telegram Finance Bot is running",
		}

		// Check database health
		if err := db.Health(); err != nil {
			health["database"] = "unhealthy"
			health["database_error"] = err.Error()
			c.JSON(503, health)
			return
		}
		health["database"] = "healthy"

		c.JSON(200, health)
	})

	// Webhook endpoint for Telegram
	router.POST("/webhook", func(c *gin.Context) {
		var update telegram.Update
		if err := c.ShouldBindJSON(&update); err != nil {
			log.Printf("Error binding update: %v", err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if err := telegramHandler.HandleUpdate(update); err != nil {
			log.Printf("Error handling update: %v", err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(200, gin.H{"status": "ok"})
	})

	// Set webhook if webhook URL is configured
	if cfg.WebhookURL != "" {
		log.Printf("Setting webhook to: %s", cfg.WebhookURL)
		if _, err := telegramClient.SetWebhook(cfg.WebhookURL); err != nil {
			log.Fatalf("Failed to set webhook: %v", err)
		}
		log.Println("Webhook set successfully")
	}

	// Start server
	port := ":" + cfg.Port
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}