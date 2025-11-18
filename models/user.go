package models

const (
	RoleFaskes  = "faskes"
	RoleAdmin   = "admin"
	RoleAuditor = "auditor"
)

type User struct {
	ID           string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string `gorm:"uniqueIndex;type:varchar(64);not null" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Role         string `gorm:"type:varchar(32);not null" json:"role"` // admin, faskes, auditor
}
