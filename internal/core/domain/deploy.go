package domain

import "time"

// DeployEvent represents a deployment event
type DeployEvent struct {
	Repository string    // ชื่อ repository
	Branch     string    // ชื่อ branch
	Commit     struct {
		SHA     string   // รหัส commit
		Message string   // ข้อความ commit
		Author  string   // ผู้ commit
	}
	Timestamp time.Time  // เวลาที่ deploy
}
