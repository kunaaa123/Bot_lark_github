package secondary

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestGitHubAdapter_ValidateWebhook(t *testing.T) {
	// สร้าง test secret และ payload
	secret := "test-secret"
	payload := []byte(`{"repository":{"full_name":"test/repo"}}`)

	// คำนวณ signature ที่ถูกต้อง
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name      string
		secret    string
		payload   []byte
		signature string
		wantErr   bool
	}{
		{
			name:      "valid signature",
			secret:    secret,
			payload:   payload,
			signature: validSignature,
			wantErr:   false,
		},
		{
			name:      "empty signature",
			secret:    secret,
			payload:   payload,
			signature: "",
			wantErr:   true,
		},
		{
			name:      "invalid signature",
			secret:    secret,
			payload:   payload,
			signature: "sha256=invalid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewGitHubAdapter(tt.secret)
			err := adapter.ValidateWebhook(tt.payload, tt.signature)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitHubAdapter_HandleDeployEvent(t *testing.T) {
	// Mock payload data
	payload := []byte(`{
        "repository": {
            "name": "test-repo",
            "full_name": "test/repo"
        },
        "head_commit": {
            "id": "abc123",
            "message": "test commit",
            "author": {
                "name": "Test User"
            },
            "timestamp": "2025-06-03T10:00:00Z"
        },
        "ref": "refs/heads/main"
    }`)

	adapter := NewGitHubAdapter("test-secret")
	event, err := adapter.HandleDeployEvent(payload)

	if err != nil {
		t.Fatalf("HandleDeployEvent() error = %v", err)
	}

	// Test each field
	if event.Repository != "test/repo" {
		t.Errorf("Expected repository 'test/repo', got '%s'", event.Repository)
	}
	if event.Branch != "main" {
		t.Errorf("Expected branch 'main', got '%s'", event.Branch)
	}
	if event.Commit.SHA != "abc123" {
		t.Errorf("Expected SHA 'abc123', got '%s'", event.Commit.SHA)
	}
}
