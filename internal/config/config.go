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
		larkWebhookURL = "https://open.larksuite.com/open-apis/bot/v2/hook/66a2d4a9-a7dd-47d3-a15a-c11c6f97c7f5"
	}

	return &Config{
		Port:           port,
		LarkWebhookURL: larkWebhookURL,
	}
}
