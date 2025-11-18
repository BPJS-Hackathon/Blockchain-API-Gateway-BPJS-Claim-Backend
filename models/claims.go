package models

import "time"

type ClaimStatus string

const (
	StatusPending   ClaimStatus = "pending"
	StatusVerifying ClaimStatus = "verifying"
	StatusApproved  ClaimStatus = "approved"
	StatusRejected  ClaimStatus = "rejected"
	StatusSubmitted ClaimStatus = "submitted"
	StatusConfirmed ClaimStatus = "confirmed"
)

type Claim struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	FacilityID  string    `gorm:"not null" json:"facility_id"`
	PatientID   string    `gorm:"not null" json:"patient_id"`
	DiagnosisID string    `gorm:"not null" json:"diagnosis_id"`
	Amount      int64     `gorm:"not null" json:"amount"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
