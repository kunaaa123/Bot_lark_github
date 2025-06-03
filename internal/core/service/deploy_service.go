package service

import (
	"bot-lark-github/internal/core/domain"
	"bot-lark-github/internal/core/ports"
	"time"
)

type DeployService struct {
	notificationService ports.NotificationService
	gitRepository       ports.GitRepository
}

func NewDeployService(notificationService ports.NotificationService, gitRepository ports.GitRepository) *DeployService {
	return &DeployService{
		notificationService: notificationService,
		gitRepository:       gitRepository,
	}
}

func (s *DeployService) ProcessDeployment(info domain.DeploymentInfo) error {
	info.Timestamp = time.Now()
	return s.notificationService.SendDeploymentNotification(info)
}

func (s *DeployService) ProcessGitWebhook(payload []byte, eventType string) error {
	if eventType != "push" {
		return nil
	}

	pushEvent, err := s.gitRepository.ParsePushEvent(payload)
	if err != nil {
		return err
	}
	commitInfo := s.gitRepository.ConvertToGitCommitInfo(pushEvent)
	return s.notificationService.SendGitDeploymentNotification(commitInfo)
}

func (s *DeployService) ProcessTestNotification() error {
	info := domain.DeploymentInfo{
		Environment: "TEST",
		Deployer:    "test-user",
		ServiceName: "test-service",
		CommitMsg:   "test: testing notification system",
		RepoURL:     "https://github.com/kunaaa123/Bot_Test", // Replace with actual repo URL
	}

	return s.notificationService.SendDeploymentNotification(info)
}

func (s *DeployService) ProcessCustomNotification(info domain.DeploymentInfo) error {
	info.Timestamp = time.Now()
	return s.notificationService.SendDeploymentNotification(info)
}
