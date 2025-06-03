package main

import (
	"log"
	"net/http"

	"bot-lark-github/internal/adapters/primary"
	"bot-lark-github/internal/adapters/secondary"
	"bot-lark-github/internal/config"
	"bot-lark-github/internal/core/service"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize adapters
	gitHubAdapter := secondary.NewGitHubAdapter(cfg.GitHubSecret)
	larkAdapter := secondary.NewLarkAdapter(cfg.LarkWebhookURL)

	// Create deploy service with dependencies
	deployService := service.NewDeployService(
		gitHubAdapter,
		larkAdapter,
	)

	// Create webhook handler
	webhookHandler := primary.NewWebhookHandler(deployService)

	// Setup routes
	http.HandleFunc("/webhook", webhookHandler.HandleWebhook)

	// Start server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
