package ports

import "bot-lark-github/internal/core/domain"

// GitPort defines interface for GitHub webhook operations
type GitPort interface {
	// ตรวจสอบความถูกต้องของ webhook
	ValidateWebhook(payload []byte, signature string) error

	// แปลงข้อมูล webhook เป็น DeployEvent
	HandleDeployEvent(payload []byte) (*domain.DeployEvent, error)
}
