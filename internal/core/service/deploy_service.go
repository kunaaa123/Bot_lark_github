package service

import (
	"bot-lark-github/internal/core/ports"
)

type DeployService struct {
	gitPort          ports.GitPort
	notificationPort ports.NotificationPort
}

func NewDeployService(
	gitPort ports.GitPort,
	notificationPort ports.NotificationPort,
) *DeployService {
	return &DeployService{
		gitPort:          gitPort,
		notificationPort: notificationPort,
	}
}

func (s *DeployService) HandleDeploy(payload []byte, signature string) error {
	// Validate webhook
	if err := s.gitPort.ValidateWebhook(payload, signature); err != nil {
		return err
	}

	// Parse deployment event
	event, err := s.gitPort.HandleDeployEvent(payload)
	if err != nil {
		return err
	}

	// Send notification
	return s.notificationPort.SendNotification(event)
}
