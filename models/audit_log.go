package models

import "time"

type AuditLog struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	UserID    *string        `gorm:"type:uuid" json:"user_id"`
	Action    string         `gorm:"type:text;not null" json:"action"`
	Meta      map[string]any `gorm:"type:jsonb;default:'{}'" json:"meta"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
}
