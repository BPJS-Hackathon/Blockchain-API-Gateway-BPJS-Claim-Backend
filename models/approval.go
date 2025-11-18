package models

import "time"

type Approval struct {
	ID         string    `gorm:"primaryKey;type:uuid" json:"id"`
	ClaimID    string    `gorm:"type:uuid;index;not null" json:"claim_id"`
	ApproverID string    `gorm:"type:uuid" json:"approver_id"`
	Status     string    `gorm:"type:text;not null" json:"status"` // approve|reject
	Note       *string   `gorm:"type:text" json:"note"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
 