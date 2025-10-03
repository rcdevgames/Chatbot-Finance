package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"chatbot/internal/config"
	"chatbot/internal/database"
	"chatbot/internal/model"
	"chatbot/internal/service"

	"github.com/joho/godotenv"

	"gorm.io/gorm"
)

func main() {
	var (
		action      = flag.String("action", "help", "License action: generate, list, activate, deactivate, status, validate, help")
		name        = flag.String("name", "", "License name (for generate)")
		description = flag.String("description", "", "License description (for generate)")
		maxUsers    = flag.Int("max-users", 1, "Maximum users allowed (0 for unlimited, for generate)")
		daysValid   = flag.Int("days", 365, "Days valid (0 for no expiry, for generate)")
		createdBy   = flag.String("created-by", "admin", "Who created this license (for generate)")
		licenseKey  = flag.String("key", "", "License key (for activate, deactivate, validate)")
		userID      = flag.String("user-id", "", "User ID for activation/deactivation")
		configPath  = flag.String("config", ".env", "Path to config file")
		force       = flag.Bool("force", false, "Force action without confirmation")
	)
	flag.Parse()

	// Load configuration
	// Load .env file from specified path if provided
	if *configPath != ".env" {
		if err := godotenv.Load(*configPath); err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	cfg := config.Load()

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize license service
	licenseService := service.NewLicenseService(db.DB)

	// Execute license action
	switch *action {
	case "generate":
		if err := generateLicense(licenseService, *name, *description, *maxUsers, *daysValid, *createdBy); err != nil {
			log.Fatalf("Failed to generate license: %v", err)
		}
	case "list":
		if err := listLicenses(db.DB); err != nil {
			log.Fatalf("Failed to list licenses: %v", err)
		}
	case "activate":
		if err := activateLicenseForUser(licenseService, *licenseKey, *userID); err != nil {
			log.Fatalf("Failed to activate license: %v", err)
		}
	case "deactivate":
		if err := deactivateLicenseForUser(licenseService, *userID, *force); err != nil {
			log.Fatalf("Failed to deactivate license: %v", err)
		}
	case "status":
		if *userID != "" {
			if err := showUserLicenseStatus(licenseService, *userID); err != nil {
				log.Fatalf("Failed to show user license status: %v", err)
			}
		} else if *licenseKey != "" {
			if err := showLicenseKeyStatus(licenseService, *licenseKey); err != nil {
				log.Fatalf("Failed to show license status: %v", err)
			}
		} else {
			if err := showAllLicenseStatus(db.DB); err != nil {
				log.Fatalf("Failed to show license status: %v", err)
			}
		}
	case "validate":
		if *licenseKey == "" {
			fmt.Println("License key is required for validation. Use -key flag.")
			os.Exit(1)
		}
		if err := validateLicense(licenseService, *licenseKey); err != nil {
			log.Fatalf("Failed to validate license: %v", err)
		}
	case "help":
		showHelp()
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		showHelp()
		os.Exit(1)
	}
}

func generateLicense(ls *service.LicenseService, name, description string, maxUsers, daysValid int, createdBy string) error {
	if name == "" {
		return fmt.Errorf("license name is required")
	}

	// Create license generator
	generator := model.NewLicenseGenerator("FBT") // FinanceBot prefix
	licenseKey := generator.GenerateKey()

	// Calculate expiry date
	var expiresAt *time.Time
	if daysValid > 0 {
		expiry := time.Now().AddDate(0, 0, daysValid)
		expiresAt = &expiry
	}

	// Create license
	license := &model.LicenseKey{
		Key:         licenseKey,
		Name:        name,
		Description: description,
		MaxUsers:    maxUsers,
		IsActive:    true,
		ExpiresAt:   expiresAt,
		CreatedBy:   createdBy,
	}

	if err := ls.CreateLicense(license); err != nil {
		return fmt.Errorf("failed to create license: %w", err)
	}

	fmt.Printf("✅ License generated successfully!\n\n")
	fmt.Printf("🔑 License Key: %s\n", licenseKey)
	fmt.Printf("📝 Name: %s\n", name)
	fmt.Printf("📄 Description: %s\n", description)
	fmt.Printf("👥 Max Users: %d\n", maxUsers)
	fmt.Printf("⏰ Valid for: %d days\n", daysValid)
	fmt.Printf("👤 Created by: %s\n", createdBy)

	if expiresAt != nil {
		fmt.Printf("📅 Expires on: %s\n", expiresAt.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("\n📋 Usage:\n")
	fmt.Printf("User can activate with: /license %s\n", licenseKey)

	return nil
}

func listLicenses(db *gorm.DB) error {
	var licenses []model.LicenseKey
	if err := db.Order("created_at DESC").Find(&licenses).Error; err != nil {
		return fmt.Errorf("failed to fetch licenses: %w", err)
	}

	if len(licenses) == 0 {
		fmt.Println("No licenses found.")
		return nil
	}

	fmt.Printf("📋 LICENSE LIST (%d total)\n\n", len(licenses))
	fmt.Printf("%-20s %-25s %-10s %-8s %-12s %-15s %s\n",
		"KEY", "NAME", "STATUS", "USERS", "MAX_USERS", "EXPIRES", "CREATED_BY")
	fmt.Println(strings.Repeat("-", 110))

	for _, license := range licenses {
		// Count active users
		var activeUsers int64
		db.Model(&model.UserLicense{}).
			Where("license_key_id = ? AND is_active = ?", license.ID, true).
			Count(&activeUsers)

		status := "Active"
		if !license.IsActive {
			status = "Inactive"
		} else if license.ExpiresAt != nil && license.ExpiresAt.Before(time.Now()) {
			status = "Expired"
		}

		expires := "Never"
		if license.ExpiresAt != nil {
			expires = license.ExpiresAt.Format("2006-01-02")
		}

		maxUsers := "Unlimited"
		if license.MaxUsers > 0 {
			maxUsers = strconv.Itoa(license.MaxUsers)
		}

		// Truncate key for display
		displayKey := license.Key
		if len(displayKey) > 20 {
			displayKey = displayKey[:17] + "..."
		}

		// Truncate name for display
		displayName := license.Name
		if len(displayName) > 25 {
			displayName = displayName[:22] + "..."
		}

		fmt.Printf("%-20s %-25s %-10s %-8d %-12s %-15s %s\n",
			displayKey, displayName, status, activeUsers, maxUsers, expires, license.CreatedBy)
	}

	return nil
}

func activateLicenseForUser(ls *service.LicenseService, licenseKey, userIDStr string) error {
	if licenseKey == "" {
		return fmt.Errorf("license key is required")
	}
	if userIDStr == "" {
		return fmt.Errorf("user ID is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Get or create user (this is a simplified version - in real implementation you'd get user from UserService)
	// For now, we'll just validate the license
	status, err := ls.ValidateLicense(licenseKey)
	if err != nil {
		return fmt.Errorf("license validation failed: %w", err)
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

	fmt.Printf("✅ License %s is valid and can be activated by user %d\n", licenseKey, userID)
	fmt.Printf("📊 License: %s\n", status.LicenseName)
	fmt.Printf("👥 Users: %d/%d\n", status.UsersActive, status.MaxUsers)

	if status.ExpiresAt != nil && status.DaysRemaining != nil {
		fmt.Printf("⏰ Expires in: %d days\n", *status.DaysRemaining)
	}

	fmt.Printf("\n📋 User can activate with: /license %s\n", licenseKey)

	return nil
}

func deactivateLicenseForUser(ls *service.LicenseService, userIDStr string, force bool) error {
	if userIDStr == "" {
		return fmt.Errorf("user ID is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if !force {
		fmt.Printf("⚠️  This will deactivate license for user %d. Are you sure? (y/N): ", userID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Deactivation cancelled")
			return nil
		}
	}

	if err := ls.DeactivateLicense(userID); err != nil {
		return fmt.Errorf("failed to deactivate license: %w", err)
	}

	fmt.Printf("✅ License deactivated for user %d\n", userID)
	return nil
}

func showUserLicenseStatus(ls *service.LicenseService, userIDStr string) error {
	if userIDStr == "" {
		return fmt.Errorf("user ID is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	status, err := ls.GetLicenseInfo(userID)
	if err != nil {
		return fmt.Errorf("failed to get user license info: %w", err)
	}

	fmt.Printf("👤 USER LICENSE STATUS - User ID: %d\n\n", userID)
	fmt.Printf("🔑 License: %s\n", status.LicenseName)
	fmt.Printf("📊 Status: ", status.LicenseName)

	if status.IsValid {
		fmt.Printf("✅ Valid\n")
	} else {
		if status.IsExpired {
			fmt.Printf("❌ Expired\n")
		} else {
			fmt.Printf("❌ Invalid\n")
		}
	}

	fmt.Printf("👥 Users: %d/%d\n", status.UsersActive, status.MaxUsers)

	if status.ExpiresAt != nil {
		fmt.Printf("📅 Expires: %s\n", status.ExpiresAt.Format("2006-01-02 15:04:05"))
		if status.DaysRemaining != nil {
			fmt.Printf("⏰ Days remaining: %d\n", *status.DaysRemaining)
		}
	}

	return nil
}

func showLicenseKeyStatus(ls *service.LicenseService, licenseKey string) error {
	if licenseKey == "" {
		return fmt.Errorf("license key is required")
	}

	status, err := ls.ValidateLicense(licenseKey)
	if err != nil {
		return fmt.Errorf("failed to validate license: %w", err)
	}

	fmt.Printf("🔑 LICENSE KEY STATUS - %s\n\n", licenseKey)
	fmt.Printf("📝 License: %s\n", status.LicenseName)
	fmt.Printf("📊 Status: ", status.LicenseName)

	if status.IsValid {
		fmt.Printf("✅ Valid\n")
	} else {
		if status.IsExpired {
			fmt.Printf("❌ Expired\n")
		} else {
			fmt.Printf("❌ Invalid\n")
		}
	}

	fmt.Printf("👥 Users: %d/%d\n", status.UsersActive, status.MaxUsers)

	if status.ExpiresAt != nil {
		fmt.Printf("📅 Expires: %s\n", status.ExpiresAt.Format("2006-01-02 15:04:05"))
		if status.DaysRemaining != nil {
			fmt.Printf("⏰ Days remaining: %d\n", *status.DaysRemaining)
		}
	}

	return nil
}

func showAllLicenseStatus(db *gorm.DB) error {
	var licenses []model.LicenseKey
	if err := db.Order("created_at DESC").Find(&licenses).Error; err != nil {
		return fmt.Errorf("failed to fetch licenses: %w", err)
	}

	if len(licenses) == 0 {
		fmt.Println("No licenses found.")
		return nil
	}

	fmt.Printf("📊 LICENSE STATUS SUMMARY\n\n")

	for _, license := range licenses {
		// Count active users
		var activeUsers int64
		db.Model(&model.UserLicense{}).
			Where("license_key_id = ? AND is_active = ?", license.ID, true).
			Count(&activeUsers)

		status := "✅ Active"
		if !license.IsActive {
			status = "❌ Inactive"
		} else if license.ExpiresAt != nil && license.ExpiresAt.Before(time.Now()) {
			status = "⏰ Expired"
		}

		fmt.Printf("🔑 %s - %s\n", license.Key, license.Name)
		fmt.Printf("📊 Status: %s\n", status)
		fmt.Printf("👥 Users: %d/%d\n", activeUsers, license.MaxUsers)

		if license.ExpiresAt != nil {
			daysRemaining := int(license.ExpiresAt.Sub(time.Now()).Hours() / 24)
			if daysRemaining > 0 {
				fmt.Printf("⏰ Expires in: %d days\n", daysRemaining)
			} else {
				fmt.Printf("⏰ Expired: %s\n", license.ExpiresAt.Format("2006-01-02"))
			}
		}

		fmt.Printf("👤 Created by: %s\n\n", license.CreatedBy)
	}

	return nil
}

func validateLicense(ls *service.LicenseService, licenseKey string) error {
	status, err := ls.ValidateLicense(licenseKey)
	if err != nil {
		return fmt.Errorf("failed to validate license: %w", err)
	}

	fmt.Printf("🔍 LICENSE VALIDATION RESULT\n\n")
	fmt.Printf("🔑 License Key: %s\n", licenseKey)
	fmt.Printf("📝 License: %s\n", status.LicenseName)
	fmt.Printf("📊 Valid: ", status.LicenseName)

	if status.IsValid {
		fmt.Printf("✅ YES\n")
	} else {
		fmt.Printf("❌ NO\n")
		if status.IsExpired {
			fmt.Printf("⚠️  Reason: License has expired\n")
		}
		if status.UsersActive >= status.MaxUsers && status.MaxUsers > 0 {
			fmt.Printf("⚠️  Reason: User limit reached (%d/%d)\n", status.UsersActive, status.MaxUsers)
		}
	}

	fmt.Printf("👥 Users: %d/%d\n", status.UsersActive, status.MaxUsers)

	if status.ExpiresAt != nil && status.DaysRemaining != nil {
		fmt.Printf("⏰ Days remaining: %d\n", *status.DaysRemaining)
	}

	return nil
}

func showHelp() {
	fmt.Println(`🔑 LICENSE MANAGEMENT CLI

USAGE:
    go run cmd/license/main.go [FLAGS]

ACTIONS:
    generate    Generate a new license key
    list        List all licenses
    activate    Activate license for user (validation only)
    deactivate  Deactivate license for user
    status      Show license status
    validate    Validate license key
    help        Show this help

FLAGS:
    -action string         Action to perform (default: "help")
    -name string           License name (for generate)
    -description string    License description (for generate)
    -max-users int         Maximum users allowed (default: 1)
    -days int              Days valid (default: 365, 0 = no expiry)
    -created-by string     Creator of license (default: "admin")
    -key string            License key (for activate, deactivate, validate)
    -user-id string        User ID for activation/deactivation
    -config string         Config file path (default: ".env")
    -force                 Force action without confirmation

EXAMPLES:
    # Generate a personal license
    go run cmd/license/main.go -action generate -name "Personal License" -max-users 1 -days 365

    # Generate a team license
    go run cmd/license/main.go -action generate -name "Team License" -max-users 10 -days 365

    # Generate unlimited license
    go run cmd/license/main.go -action generate -name "Enterprise License" -max-users 0 -days 0

    # List all licenses
    go run cmd/license/main.go -action list

    # Validate license key
    go run cmd/license/main.go -action validate -key "FBT-1234-5678-9ABC-DEF0"

    # Check license status
    go run cmd/license/main.go -action status -key "FBT-1234-5678-9ABC-DEF0"

    # Check user license status
    go run cmd/license/main.go -action status -user-id "123456789"

    # Deactivate user license
    go run cmd/license/main.go -action deactivate -user-id "123456789" -force`)
}