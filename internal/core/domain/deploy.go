package domain

import "time"

type DeploymentInfo struct {
	Environment string
	Deployer    string
	ServiceName string
	CommitMsg   string
	RepoURL     string
	Timestamp   time.Time
}

type GitCommitInfo struct {
	Message     string
	Features    []string
	Environment string
	ServiceName string
	Deployer    string
	RepoURL     string
}

type GitHubPushEvent struct {
	Ref        string
	Repository Repository
	Commits    []Commit
}

type Repository struct {
	Name string
}

type Commit struct {
	Message string
	Author  Author
}

type Author struct {
	Name  string
	Email string
}

type NotificationCard struct {
	Title       string
	Template    string
	Environment string
	Deployer    string
	ServiceName string
	Message     string
	RepoURL     string
	Timestamp   time.Time
	Actions     []Action
}

////
type Action struct {
	Text string
	URL  string
	Type string
}
