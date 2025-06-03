package ports

import "bot-lark-github/internal/core/domain"

type NotificationService interface {
	SendDeploymentNotification(info domain.DeploymentInfo) error
	SendGitDeploymentNotification(info domain.GitCommitInfo) error
	BuildNotificationCard(info domain.DeploymentInfo) domain.NotificationCard
	BuildGitNotificationCard(info domain.GitCommitInfo) domain.NotificationCard
}