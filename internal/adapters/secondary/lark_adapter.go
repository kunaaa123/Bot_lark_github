package secondary

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"bot-lark-github/internal/core/domain"
)

// LarkAdapter implements NotificationPort interface
type LarkAdapter struct {
	webhookURL string
	client     *http.Client
}

// NewLarkAdapter สร้าง instance ใหม่ของ LarkAdapter
func NewLarkAdapter(webhookURL string) *LarkAdapter {
	return &LarkAdapter{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// โครงสร้างข้อมูลสำหรับส่งไป Lark
type larkMessage struct {
	MsgType string `json:"msg_type"`
	Card    struct {
		Header struct {
			Title struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"title"`
		} `json:"header"`
		Elements []struct {
			Tag     string `json:"tag"`
			Text    string `json:"text,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"elements"`
	} `json:"card"`
}

// custom errors
var (
	ErrEmptyWebhookURL = errors.New("webhook URL is empty")
	ErrSendMessage     = errors.New("failed to send message to Lark")
)

// SendNotification ส่งการแจ้งเตือนไปยัง Lark
func (a *LarkAdapter) SendNotification(event *domain.DeployEvent) error {
	if a.webhookURL == "" {
		return ErrEmptyWebhookURL
	}

	// สร้าง message card
	msg := &larkMessage{
		MsgType: "interactive",
	}

	// set header
	msg.Card.Header.Title.Tag = "plain_text"
	msg.Card.Header.Title.Content = fmt.Sprintf("🚀 New Deployment: %s", event.Repository)

	// add elements
	msg.Card.Elements = []struct {
		Tag     string `json:"tag"`
		Text    string `json:"text,omitempty"`
		Content string `json:"content,omitempty"`
	}{
		{
			Tag: "markdown",
			Text: fmt.Sprintf("**Branch:** %s\n**Commit:** %s\n**Author:** %s\n**Message:** %s",
				event.Branch,
				event.Commit.SHA,
				event.Commit.Author,
				event.Commit.Message),
		},
		{
			Tag:     "note",
			Content: fmt.Sprintf("Deployed at: %s", event.Timestamp.Format(time.RFC3339)),
		},
	}

	// แปลงเป็น JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// ส่ง HTTP POST request
	resp, err := a.client.Post(
		a.webhookURL,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSendMessage, err)
	}
	defer resp.Body.Close()

	// ตรวจสอบ response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status code %d", ErrSendMessage, resp.StatusCode)
	}

	return nil
}
