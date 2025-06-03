package primary

import (
	"bot-lark-github/internal/core/domain"
	"bot-lark-github/internal/core/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HTTPHandler struct {
	deployService *service.DeployService
}

func NewHTTPHandler(deployService *service.DeployService) *HTTPHandler {
	return &HTTPHandler{
		deployService: deployService,
	}
}

func (h *HTTPHandler) HandleDeploymentInfo(w http.ResponseWriter, r *http.Request) {
	info := domain.DeploymentInfo{
		Environment: "DEV",
		Deployer:    "rutchanai",
		ServiceName: "tgth-backend-main",
		CommitMsg:   "feat: backend deployment structure",
		RepoURL:     "https://github.com/kunaaa123/Bot_Test",
	}

	if err := h.deployService.ProcessDeployment(info); err != nil {
		log.Printf("Failed to process deployment: %v", err)
		http.Error(w, "Failed to process deployment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (h *HTTPHandler) HandleTestNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.deployService.ProcessTestNotification(); err != nil {
		log.Printf("Error processing test notification: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "notification sent successfully",
	})
}

func (h *HTTPHandler) HandleCustomNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var info domain.DeploymentInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.deployService.ProcessCustomNotification(info); err != nil {
		log.Printf("Error processing custom notification: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "custom notification sent successfully",
	})
}

func (h *HTTPHandler) HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received webhook from GitHub")

	eventType := r.Header.Get("X-GitHub-Event")
	
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := h.deployService.ProcessGitWebhook(payload, eventType); err != nil {
		log.Printf("Error processing webhook: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook processed successfully"))
}