package ports

import "bot-lark-github/internal/core/domain"

type GitRepository interface {
	ParsePushEvent(payload []byte) (*domain.GitHubPushEvent, error)
	ExtractCommitMessages(event *domain.GitHubPushEvent) string
	ConvertToGitCommitInfo(event *domain.GitHubPushEvent) domain.GitCommitInfo
}