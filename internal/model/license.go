package model

import (
	"time"

	"github.com/google/uuid"
)

// LicenseKey represents a license key for accessing the bot
type LicenseKey struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Key         string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"key"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	MaxUsers    int       `gorm:"default:1" json:"max_users"`        // Maximum users allowed
	IsActive    bool      `gorm:"default:true" json:"is_active"`     // Is license active
	ExpiresAt   *time.Time `gorm:"index" json:"expires_at"`          // License expiry date
	CreatedBy   string    `gorm:"type:varchar(100)" json:"created_by"` // Who created this license
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserLicense represents a user activated with a license
type UserLicense struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	LicenseKeyID string    `gorm:"type:uuid;not null;index" json:"license_key_id"`
	UserID       string    `gorm:"type:uuid;not null;index" json:"user_id"`
	ActivatedAt  time.Time `gorm:"autoCreateTime" json:"activated_at"`
	ActivatedBy  int64     `gorm:"not null;index" json:"activated_by"` // Telegram user ID who activated
	IsActive     bool      `gorm:"default:true" json:"is_active"`     // Is this activation active
	LastUsed     *time.Time `json:"last_used"`                        // Last time this user used the bot
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	LicenseKey LicenseKey `gorm:"foreignKey:LicenseKeyID"`
	User       User       `gorm:"foreignKey:UserID"`
}

func (LicenseKey) TableName() string {
	return "license_keys"
}

func (UserLicense) TableName() string {
	return "user_licenses"
}

// LicenseStatus represents the status of a license
type LicenseStatus struct {
	IsValid       bool      `json:"is_valid"`
	IsExpired     bool      `json:"is_expired"`
	UsersActive   int       `json:"users_active"`
	MaxUsers      int       `json:"max_users"`
	ExpiresAt     *time.Time `json:"expires_at"`
	DaysRemaining *int      `json:"days_remaining,omitempty"`
	LicenseName   string    `json:"license_name"`
}

// LicenseGenerator generates new license keys
type LicenseGenerator struct {
	Prefix string // e.g., "FBT" for FinanceBot
}

// NewLicenseGenerator creates a new license generator
func NewLicenseGenerator(prefix string) *LicenseGenerator {
	return &LicenseGenerator{Prefix: prefix}
}

// GenerateKey generates a new license key
func (lg *LicenseGenerator) GenerateKey() string {
	// Generate UUID and format as license key
	id := uuid.New()
	// Format: PREFIX-XXXX-XXXX-XXXX-XXXX
	key := id.String()
	if lg.Prefix != "" {
		key = lg.Prefix + "-" + key[:8] + "-" + key[9:13] + "-" + key[14:18] + "-" + key[19:23]
	}
	return key
}

// Default license key configurations
var DefaultLicenseConfigs = []struct {
	Name        string
	Description string
	MaxUsers    int
	DaysValid   int // 0 for no expiry
	CreatedBy   string
}{
	{
		Name:        "Personal License",
		Description: "License for personal use - 1 user",
		MaxUsers:    1,
		DaysValid:   365, // 1 year
		CreatedBy:   "system",
	},
	{
		Name:        "Family License",
		Description: "License for family use - up to 5 users",
		MaxUsers:    5,
		DaysValid:   365,
		CreatedBy:   "system",
	},
	{
		Name:        "Team License",
		Description: "License for team use - up to 10 users",
		MaxUsers:    10,
		DaysValid:   365,
		CreatedBy:   "system",
	},
	{
		Name:        "Enterprise License",
		Description: "License for enterprise use - unlimited users",
		MaxUsers:    0, // 0 means unlimited
		DaysValid:   0,  // 0 means no expiry
		CreatedBy:   "system",
	},
}