package config

import "os"

type Config struct {
	Port           string
	LarkWebhookURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	larkWebhookURL := os.Getenv("LARK_WEBHOOK_URL")
	if larkWebhookURL == "" {
		larkWebhookURL = "https://open.larksuite.com/open-apis/bot/v2/hook/88fccfea-8fad-47d9-99a9-44d214785fff"
	}

	return &Config{
		Port:           port,
		LarkWebhookURL: larkWebhookURL,
	}
}
