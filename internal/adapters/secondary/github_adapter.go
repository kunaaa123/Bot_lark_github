package secondary

import (
	"bot-lark-github/internal/core/domain"
	"encoding/json"
)

type GitHubAdapter struct{}

func NewGitHubAdapter() *GitHubAdapter {
	return &GitHubAdapter{}
}

func (g *GitHubAdapter) ParsePushEvent(payload []byte) (*domain.GitHubPushEvent, error) {
	var pushEvent domain.GitHubPushEvent
	if err := json.Unmarshal(payload, &pushEvent); err != nil {
		return nil, err
	}
	return &pushEvent, nil
}

func (g *GitHubAdapter) ExtractCommitMessages(event *domain.GitHubPushEvent) string {
	commitMessages := ""
	for _, commit := range event.Commits {
		commitMessages += "- " + commit.Message + "\n"
	}
	return commitMessages
}

func (g *GitHubAdapter) ConvertToGitCommitInfo(event *domain.GitHubPushEvent) domain.GitCommitInfo {
	commitMessages := g.ExtractCommitMessages(event)
	
	deployer := "Unknown"
	if len(event.Commits) > 0 {
		deployer = event.Commits[0].Author.Name
	}

	return domain.GitCommitInfo{
		Message:     commitMessages,
		Environment: "DEV",
		ServiceName: event.Repository.Name,
		Deployer:    deployer,
		RepoURL:     "https://github.com/kunaaa123/Bot_Test",
	}
}