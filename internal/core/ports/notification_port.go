package ports

import "bot-lark-github/internal/core/domain"

// NotificationPort defines interface for sending notifications
type NotificationPort interface {
	// ส่งการแจ้งเตือนไปยัง Lark
	SendNotification(event *domain.DeployEvent) error
}
