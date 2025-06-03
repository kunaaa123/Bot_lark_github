package primary

import (
	"io/ioutil"
	"net/http"

	"bot-lark-github/internal/core/service"
)

// WebhookHandler handles incoming GitHub webhooks
type WebhookHandler struct {
	deployService *service.DeployService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(deployService *service.DeployService) *WebhookHandler {
	return &WebhookHandler{
		deployService: deployService,
	}
}

// HandleWebhook processes GitHub webhook requests
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบ method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// อ่าน signature จาก header
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		http.Error(w, "Missing signature", http.StatusBadRequest)
		return
	}

	// อ่าน payload
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// ส่งต่อไปให้ Deploy Service จัดการ
	err = h.deployService.HandleDeploy(payload, signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
