package secondary

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bot-lark-github/internal/core/domain"
)

func TestLarkAdapter_SendNotification(t *testing.T) {
	// สร้าง test data
	testEvent := &domain.DeployEvent{
		Repository: "test/repo",
		Branch:     "main",
		Commit: struct {
			SHA     string
			Message string
			Author  string
		}{
			SHA:     "abc123",
			Message: "test commit",
			Author:  "test user",
		},
		Timestamp: time.Now(),
	}

	tests := []struct {
		name          string
		webhookURL    string
		event         *domain.DeployEvent
		serverHandler func(w http.ResponseWriter, r *http.Request)
		wantErr       bool
	}{
		{
			name:       "successful notification",
			webhookURL: "http://test.com",
			event:      testEvent,
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				// ตรวจสอบ request method
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// ตรวจสอบ content type
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected application/json, got %s", r.Header.Get("Content-Type"))
				}

				// อ่าน request body
				body, _ := ioutil.ReadAll(r.Body)
				var msg larkMessage
				if err := json.Unmarshal(body, &msg); err != nil {
					t.Errorf("Failed to parse request body: %v", err)
				}

				// ตรวจสอบข้อมูลใน message
				if msg.MsgType != "interactive" {
					t.Errorf("Expected interactive message, got %s", msg.MsgType)
				}

				w.WriteHeader(http.StatusOK)
			},
			wantErr: false,
		},
		{
			name:       "empty webhook URL",
			webhookURL: "",
			event:      testEvent,
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true,
		},
		{
			name:       "server error",
			webhookURL: "http://test.com",
			event:      testEvent,
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// สร้าง test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverHandler))
			defer server.Close()

			// สร้าง adapter โดยใช้ URL ของ test server
			adapter := NewLarkAdapter(server.URL)
			if tt.webhookURL == "" {
				adapter = NewLarkAdapter("")
			}

			// ทดสอบการส่ง notification
			err := adapter.SendNotification(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
