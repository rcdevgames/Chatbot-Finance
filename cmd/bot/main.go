package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"chatbot/internal/config"
	"chatbot/internal/handler"
	"chatbot/internal/repository"
	"chatbot/internal/service"
	"chatbot/pkg/telegram"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize dependencies
	telegramClient := telegram.NewClient(cfg.TelegramBotToken)
	supabaseClient := repository.NewSupabaseClient(cfg.SupabaseURL, cfg.SupabaseKey)
	groqService := service.NewGroqService(cfg.GroqAPIKey)

	// Initialize repositories
	userRepo := repository.NewUserRepository(supabaseClient)
	transactionRepo := repository.NewTransactionRepository(supabaseClient)
	chatHistoryRepo := repository.NewChatHistoryRepository(supabaseClient)
	categoryRepo := repository.NewCategoryRepository(supabaseClient)

	// Initialize services
	userService := service.NewUserService(userRepo, transactionRepo, chatHistoryRepo, categoryRepo)

	// Initialize handlers
	telegramHandler := handler.NewTelegramHandler(telegramClient, userService, groqService)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Telegram Finance Bot is running",
		})
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