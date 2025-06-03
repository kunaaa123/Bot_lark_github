package service

import (
	"bot-lark-github/internal/core/ports"
)

type DeployService struct {
	GitPort ports.GitPort
	NotificationPort ports.NotificationPort
}

func NewDeployService(gitPort ports.GitPort, notificationPort ports.NotificationPort) *DeployService {
	return &DeployService{
		GitPort: gitPort,
		NotificationPort: notificationPort,
	}
}

//func (ds *DeployService) P