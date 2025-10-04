package service

import (
	"fmt"
	"log"
	"time"

	"chatbot/internal/model"

	"gorm.io/gorm"
)

// LicenseService handles license validation and management
type LicenseService struct {
	db *gorm.DB
}

// NewLicenseService creates a new license service
func NewLicenseService(db *gorm.DB) *LicenseService {
	return &LicenseService{db: db}
}

// ValidateLicense checks if a license key is valid
func (ls *LicenseService) ValidateLicense(licenseKey string) (*model.LicenseStatus, error) {
	var license model.LicenseKey
	if err := ls.db.Where("key = ? AND is_active = ?", licenseKey, true).First(&license).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &model.LicenseStatus{
				IsValid:     false,
				IsExpired:   false,
				LicenseName: "Invalid License",
			}, fmt.Errorf("license key not found or inactive")
		}
		return nil, fmt.Errorf("failed to validate license: %w", err)
	}

	// Check if license is expired
	isExpired := false
	if license.ExpiresAt != nil && license.ExpiresAt.Before(time.Now()) {
		isExpired = true
	}

	// Count active users for this license
	var activeUsers int64
	ls.db.Model(&model.UserLicense{}).
		Where("license_key_id = ? AND is_active = ?", license.ID, true).
		Count(&activeUsers)

	status := &model.LicenseStatus{
		IsValid:     !isExpired,
		IsExpired:   isExpired,
		UsersActive: int(activeUsers),
		MaxUsers:    license.MaxUsers,
		ExpiresAt:   license.ExpiresAt,
		LicenseName: license.Name,
	}

	// Calculate days remaining
	if license.ExpiresAt != nil {
		daysRemaining := int(license.ExpiresAt.Sub(time.Now()).Hours() / 24)
		status.DaysRemaining = &daysRemaining
	}

	// Check if license has reached user limit
	if license.MaxUsers > 0 && int(activeUsers) >= license.MaxUsers {
		status.IsValid = false
	}

	return status, nil
}

// ActivateLicense activates a license for a user
func (ls *LicenseService) ActivateLicense(licenseKey string, telegramUserID int64, user *model.User) error {
	// Validate license first
	status, err := ls.ValidateLicense(licenseKey)
	if err != nil {
		return err
	}

	if !status.IsValid {
		if status.IsExpired {
			return fmt.Errorf("license has expired")
		}
		if status.UsersActive >= status.MaxUsers && status.MaxUsers > 0 {
			return fmt.Errorf("license has reached maximum user limit (%d users)", status.MaxUsers)
		}
		return fmt.Errorf("license is not valid")
	}

	// Get the license
	var license model.LicenseKey
	if err := ls.db.Where("key = ?", licenseKey).First(&license).Error; err != nil {
		return fmt.Errorf("failed to get license: %w", err)
	}

	// Check if user already has this license activated
	var existingUserLicense model.UserLicense
	if err := ls.db.Joins("JOIN users ON user_licenses.user_id = users.id").
		Where("user_licenses.license_key_id = ? AND users.telegram_user_id = ?", license.ID, telegramUserID).
		First(&existingUserLicense).Error; err == nil {
		if existingUserLicense.IsActive {
			return fmt.Errorf("license already activated for this user")
		}
		// Reactivate existing license
		existingUserLicense.IsActive = true
		existingUserLicense.ActivatedAt = time.Now()
		if err := ls.db.Save(&existingUserLicense).Error; err != nil {
			return fmt.Errorf("failed to reactivate license: %w", err)
		}
		log.Printf("Reactivated license %s for user %d", licenseKey, telegramUserID)
		return nil
	}

	// Create new user license activation
	userLicense := model.UserLicense{
		LicenseKeyID: license.ID,
		UserID:       user.ID,
		ActivatedBy:  telegramUserID,
		IsActive:     true,
	}

	if err := ls.db.Create(&userLicense).Error; err != nil {
		return fmt.Errorf("failed to activate license: %w", err)
	}

	log.Printf("Successfully activated license %s for user %d", licenseKey, telegramUserID)
	return nil
}

// DeactivateLicense deactivates a license for a user
func (ls *LicenseService) DeactivateLicense(telegramUserID int64) error {
	var userLicense model.UserLicense
	if err := ls.db.Joins("JOIN users ON user_licenses.user_id = users.id").
		Where("users.telegram_user_id = ? AND user_licenses.is_active = ?", telegramUserID, true).
		First(&userLicense).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no active license found for user %d", telegramUserID)
		}
		return fmt.Errorf("failed to find user license: %w", err)
	}

	userLicense.IsActive = false
	if err := ls.db.Save(&userLicense).Error; err != nil {
		return fmt.Errorf("failed to deactivate license: %w", err)
	}

	log.Printf("Deactivated license for user %d", telegramUserID)
	return nil
}

// IsUserLicensed checks if a user has a valid active license
func (ls *LicenseService) IsUserLicensed(telegramUserID int64) (bool, error) {
	var count int64
	err := ls.db.Model(&model.UserLicense{}).
		Joins("JOIN users ON user_licenses.user_id = users.id").
		Joins("JOIN license_keys ON user_licenses.license_key_id = license_keys.id").
		Where("users.telegram_user_id = ? AND user_licenses.is_active = ? AND license_keys.is_active = ?",
			telegramUserID, true, true).
		Where("(license_keys.expires_at IS NULL OR license_keys.expires_at > ?)", time.Now()).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check user license: %w", err)
	}

	return count > 0, nil
}

// UpdateLastUsed updates the last used timestamp for a user license
func (ls *LicenseService) UpdateLastUsed(telegramUserID int64) error {
	now := time.Now()

	// First get the user ID from telegram user ID
	var user model.User
	if err := ls.db.Where("telegram_user_id = ?", telegramUserID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Update the user license using the user ID
	err := ls.db.Model(&model.UserLicense{}).
		Where("user_id = ? AND is_active = ?", user.ID, true).
		Update("last_used", now).Error

	if err != nil {
		return fmt.Errorf("failed to update last used timestamp: %w", err)
	}

	return nil
}

// GetLicenseInfo gets license information for a user
func (ls *LicenseService) GetLicenseInfo(telegramUserID int64) (*model.LicenseStatus, error) {
	var license model.LicenseKey
	err := ls.db.Joins("JOIN user_licenses ON license_keys.id = user_licenses.license_key_id").
		Joins("JOIN users ON user_licenses.user_id = users.id").
		Where("users.telegram_user_id = ? AND user_licenses.is_active = ?", telegramUserID, true).
		First(&license).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &model.LicenseStatus{
				IsValid:     false,
				LicenseName: "No License",
			}, nil
		}
		return nil, fmt.Errorf("failed to get license info: %w", err)
	}

	// Count active users for this license
	var activeUsers int64
	ls.db.Model(&model.UserLicense{}).
		Where("license_key_id = ? AND is_active = ?", license.ID, true).
		Count(&activeUsers)

	status := &model.LicenseStatus{
		IsValid:     true,
		IsExpired:   false,
		UsersActive: int(activeUsers),
		MaxUsers:    license.MaxUsers,
		ExpiresAt:   license.ExpiresAt,
		LicenseName: license.Name,
	}

	// Check if license is expired
	if license.ExpiresAt != nil && license.ExpiresAt.Before(time.Now()) {
		status.IsExpired = true
		status.IsValid = false
	}

	// Calculate days remaining
	if license.ExpiresAt != nil {
		daysRemaining := int(license.ExpiresAt.Sub(time.Now()).Hours() / 24)
		status.DaysRemaining = &daysRemaining
	}

	return status, nil
}

// GetUserLicenseKey gets the actual license key for a user
func (ls *LicenseService) GetUserLicenseKey(telegramUserID int64) (string, error) {
	var license model.LicenseKey
	err := ls.db.Joins("JOIN user_licenses ON license_keys.id = user_licenses.license_key_id").
		Joins("JOIN users ON user_licenses.user_id = users.id").
		Where("users.telegram_user_id = ? AND user_licenses.is_active = ?", telegramUserID, true).
		First(&license).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("no active license found for user")
		}
		return "", fmt.Errorf("failed to get user license key: %w", err)
	}

	return license.Key, nil
}

// CreateLicense creates a new license key
func (ls *LicenseService) CreateLicense(license *model.LicenseKey) error {
	return ls.db.Create(license).Error
}