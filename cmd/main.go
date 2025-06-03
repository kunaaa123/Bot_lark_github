package main

import (
	"bot-lark-github/internal/adapters/primary"
	"bot-lark-github/internal/adapters/secondary"
	"bot-lark-github/internal/config"
	"bot-lark-github/internal/core/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	larkAdapter := secondary.NewLarkAdapter(cfg.LarkWebhookURL)
	githubAdapter := secondary.NewGitHubAdapter()
	deployService := service.NewDeployService(larkAdapter, githubAdapter)
	httpHandler := primary.NewHTTPHandler(deployService)

	http.HandleFunc("/deployment-info", httpHandler.HandleDeploymentInfo)
	http.HandleFunc("/test-notification", httpHandler.HandleTestNotification)
	http.HandleFunc("/custom-notification", httpHandler.HandleCustomNotification)
	http.HandleFunc("/git-webhook", httpHandler.HandleGitHubWebhook)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Printf("Server running on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

