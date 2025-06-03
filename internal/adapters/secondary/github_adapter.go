package secondary

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"bot-lark-github/internal/core/domain"
)

// GitHubAdapter implements GitPort interface
type GitHubAdapter struct {
	webhookSecret string
}

// NewGitHubAdapter สร้าง instance ใหม่ของ GitHubAdapter
func NewGitHubAdapter(secret string) *GitHubAdapter {
	return &GitHubAdapter{
		webhookSecret: secret,
	}
}

// โครงสร้างข้อมูลที่รับจาก GitHub webhook เฉพาะส่วนที่ต้องการ
type githubPayload struct {
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
	HeadCommit struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name string `json:"name"`
		} `json:"author"`
		Timestamp string `json:"timestamp"`
	} `json:"head_commit"`
	Ref string `json:"ref"` // refs/heads/main
}

// เพิ่ม custom errors
var (
	ErrEmptySignature    = errors.New("signature is empty")
	ErrInvalidSignature  = errors.New("invalid signature")
	ErrEmptyPayload      = errors.New("payload is empty")
	ErrInvalidPayload    = errors.New("invalid payload format")
	ErrMissingRepository = errors.New("repository information is missing")
	ErrMissingCommit     = errors.New("commit information is missing")
)

// ValidateWebhook ตรวจสอบความถูกต้องของ webhook signature
func (a *GitHubAdapter) ValidateWebhook(payload []byte, signature string) error {
	if len(payload) == 0 {
		return ErrEmptyPayload
	}
	if signature == "" {
		return ErrEmptySignature
	}

	signature = strings.TrimPrefix(signature, "sha256=")
	mac := hmac.New(sha256.New, []byte(a.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return ErrInvalidSignature
	}

	return nil
}

// HandleDeployEvent แปลงข้อมูล webhook เป็น DeployEvent
func (a *GitHubAdapter) HandleDeployEvent(payload []byte) (*domain.DeployEvent, error) {
	if len(payload) == 0 {
		return nil, ErrEmptyPayload
	}

	var ghPayload githubPayload
	if err := json.Unmarshal(payload, &ghPayload); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}

	// ตรวจสอบข้อมูลที่จำเป็น
	if ghPayload.Repository.FullName == "" {
		return nil, ErrMissingRepository
	}
	if ghPayload.HeadCommit.ID == "" {
		return nil, ErrMissingCommit
	}

	// แปลง branch จาก ref (refs/heads/main -> main)
	branch := strings.TrimPrefix(ghPayload.Ref, "refs/heads/")

	// สร้าง DeployEvent ตาม domain model
	deployEvent := &domain.DeployEvent{
		Repository: ghPayload.Repository.FullName,
		Branch:     branch,
		Commit: struct {
			SHA     string
			Message string
			Author  string
		}{
			SHA:     ghPayload.HeadCommit.ID,
			Message: ghPayload.HeadCommit.Message,
			Author:  ghPayload.HeadCommit.Author.Name,
		},
		Timestamp: time.Now(), // ใช้เวลาปัจจุบัน
	}

	return deployEvent, nil
}
