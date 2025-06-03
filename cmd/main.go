package main

import (
	"bot-lark-github/internal/adapters/secondary"
	"bot-lark-github/internal/config"
	"fmt"
)

func main() {
	cfg := config.Load()

	LarkAdapter := secondary.NewLarkAdapter()
	GithubAdapter := secondary.NewGithubAdapter()

	deployService := service.NewDeployService(GithubAdapter, LarkAdapter)
	httpHandler := primary.NewHTTPHandler(deployService)

	

	fmt.Println("Server is running on port ...")
}