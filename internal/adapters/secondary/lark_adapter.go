package secondary

import (
	"bot-lark-github/internal/core/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LarkAdapter struct {
	webhookURL string
}

func NewLarkAdapter(webhookURL string) *LarkAdapter {
	return &LarkAdapter{
		webhookURL: webhookURL,
	}
}

func (l *LarkAdapter) SendDeploymentNotification(info domain.DeploymentInfo) error {
	card := l.BuildNotificationCard(info)
	return l.sendCard(card)
}

func (l *LarkAdapter) SendGitDeploymentNotification(info domain.GitCommitInfo) error {
	card := l.BuildGitNotificationCard(info)
	return l.sendCard(card)
}

func (l *LarkAdapter) BuildNotificationCard(info domain.DeploymentInfo) domain.NotificationCard {
	return domain.NotificationCard{
		Title:       "Backend Deployment Status",
		Template:    "indigo",
		Environment: info.Environment,
		Deployer:    info.Deployer,
		ServiceName: info.ServiceName,
		Message:     "Latest Changes:\n" + info.CommitMsg,
		RepoURL:     info.RepoURL,
		Timestamp:   info.Timestamp,
		Actions: []domain.Action{
			{
				Text: "View Repository",
				URL:  info.RepoURL,
				Type: "primary",
			},
			{
				Text: "View Documentation",
				URL:  "https://github.com/kunaaa123/Bot_Test/wiki",
				Type: "default",
			},
		},
	}
}

func (l *LarkAdapter) BuildGitNotificationCard(info domain.GitCommitInfo) domain.NotificationCard {
	return domain.NotificationCard{
		Title:       "Backend Deployment",
		Template:    "blue",
		Environment: info.Environment,
		Deployer:    info.Deployer,
		ServiceName: info.ServiceName,
		Message:     fmt.Sprintf("Commit Messages:\n%s", info.Message),
		RepoURL:     info.RepoURL,
		Timestamp:   time.Now(),
		Actions: []domain.Action{
			{
				Text: "View Repository",
				URL:  info.RepoURL,
				Type: "primary",
			},
		},
	}
}

func (l *LarkAdapter) sendCard(card domain.NotificationCard) error {
	payload := l.buildLarkPayload(card)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	log.Printf("Sending payload to Lark: %s", string(payloadBytes))

	resp, err := http.Post(l.webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error sending to Lark: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	log.Printf("Lark Response: %s", string(body))
	return nil
}

func (l *LarkAdapter) buildLarkPayload(card domain.NotificationCard) map[string]interface{} {
	elements := []map[string]interface{}{
		{
			"tag": "div",
			"text": map[string]interface{}{
				"content": fmt.Sprintf("Environment: %s\nDeployer: %s\nService: %s",
					card.Environment, card.Deployer, card.ServiceName),
				"tag": "lark_md",
			},
		},
	}

	if card.Message != "" {
		elements = append(elements, map[string]interface{}{
			"tag": "hr",
		})
		elements = append(elements, map[string]interface{}{
			"tag": "div",
			"text": map[string]interface{}{
				"content": card.Message,
				"tag":     "lark_md",
			},
		})
	}

	elements = append(elements, map[string]interface{}{
		"tag": "hr",
	})

	elements = append(elements, map[string]interface{}{
		"tag": "note",
		"elements": []map[string]interface{}{
			{
				"tag":     "plain_text",
				"content": fmt.Sprintf("Deployed at: %s", card.Timestamp.Format("2006-01-02 15:04:05")),
			},
		},
	})

	if len(card.Actions) > 0 {
		actions := make([]map[string]interface{}, len(card.Actions))
		for i, action := range card.Actions {
			actions[i] = map[string]interface{}{
				"tag": "button",
				"text": map[string]interface{}{
					"content": action.Text,
					"tag":     "plain_text",
				},
				"type": action.Type,
				"url":  action.URL,
			}
		}

		elements = append(elements, map[string]interface{}{
			"tag":     "action",
			"actions": actions,
		})
	}

	return map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]interface{}{
					"content": card.Title,
					"tag":     "plain_text",
				},
				"template": card.Template,
			},
			"elements": elements,
		},
	}
}