package database

import (
	"fmt"
	"log"
	"time"

	"chatbot/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(databaseURL string) (*Database, error) {
	log.Println("Connecting to database...")

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to database")

	// Auto migrate tables
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &Database{DB: db}, nil
}

func autoMigrate(db *gorm.DB) error {
	log.Println("Running auto migration...")

	// First, create migration schema table if not exists
	if err := db.AutoMigrate(&model.MigrationSchema{}); err != nil {
		return fmt.Errorf("failed to migrate migration schema: %w", err)
	}

	// Get current migration version
	currentVersion := getCurrentMigrationVersion(db)
	log.Printf("Current migration version: %s", currentVersion)

	// Define all model migrations
	migrations := []struct {
		version string
		name    string
		models  interface{}
	}{
		{
			version: "1.0.0",
			name:    "Initial schema with basic tables",
			models: []interface{}{
				&model.User{},
				&model.Category{},
				&model.Transaction{},
				&model.ChatHistory{},
				&model.Budget{},
				&model.Insight{},
			},
		},
		{
			version: "1.1.0",
			name:    "Add license system tables",
			models: []interface{}{
				&model.LicenseKey{},
				&model.UserLicense{},
			},
		},
	}

	// Run migrations in order
	for _, migration := range migrations {
		if shouldRunMigration(db, migration.version, currentVersion) {
			log.Printf("Running migration %s: %s", migration.version, migration.name)

			// Run AutoMigrate for models in this version
			if err := db.AutoMigrate(migration.models.([]interface{})...); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.version, err)
			}

			// Mark migration as applied
			migrationSchema := model.MigrationSchema{
				Version: migration.version,
				Name:    migration.name,
				Applied: true,
			}
			if err := db.Create(&migrationSchema).Error; err != nil {
				return fmt.Errorf("failed to mark migration %s as applied: %w", migration.version, err)
			}

			log.Printf("Migration %s completed successfully", migration.version)
		}
	}

	// Seed data after migrations
	if err := seedData(db); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	log.Println("Auto migration completed successfully")
	return nil
}

func getCurrentMigrationVersion(db *gorm.DB) string {
	var version string
	db.Model(&model.MigrationSchema{}).
		Where("applied = ?", true).
		Order("version DESC").
		Limit(1).
		Pluck("version", &version)

	if version == "" {
		return "0.0.0"
	}
	return version
}

func shouldRunMigration(db *gorm.DB, version, currentVersion string) bool {
	if version <= currentVersion {
		return false
	}

	var count int64
	db.Model(&model.MigrationSchema{}).
		Where("version = ? AND applied = ?", version, true).
		Count(&count)

	return count == 0
}

func seedData(db *gorm.DB) error {
	log.Println("Checking and seeding default data...")

	// Seed default categories if none exist
	var categoryCount int64
	db.Model(&model.Category{}).Count(&categoryCount)

	if categoryCount == 0 {
		log.Println("Seeding default categories...")
		for _, category := range model.DefaultCategories {
			if err := db.Create(&category).Error; err != nil {
				log.Printf("Failed to create category %s: %v", category.Name, err)
			}
		}
		log.Printf("Seeded %d default categories", len(model.DefaultCategories))
	} else {
		log.Printf("Categories already exist (%d records), skipping seeding", categoryCount)
	}

	// Create indexes for better performance (if not auto-created by GORM)
	if err := createAdditionalIndexes(db); err != nil {
		log.Printf("Warning: Failed to create additional indexes: %v", err)
	}

	// Create license indexes
	if err := createLicenseIndexes(db); err != nil {
		log.Printf("Warning: Failed to create license indexes: %v", err)
	}

	return nil
}

func createAdditionalIndexes(db *gorm.DB) error {
	log.Println("Creating additional indexes for performance...")

	// Composite indexes for common queries
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_transactions_user_date_type ON transactions(user_id, date DESC, type)",
		"CREATE INDEX IF NOT EXISTS idx_transactions_user_category_date ON transactions(user_id, category, date DESC)",
		"CREATE INDEX IF NOT EXISTS idx_chat_history_user_created_desc ON chat_history(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_budgets_user_active ON budgets(user_id, is_active) WHERE is_active = true",
		"CREATE INDEX IF NOT EXISTS idx_insights_user_type_period ON insights(user_id, type, period_start)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			log.Printf("Failed to create index: %s, Error: %v", indexSQL, err)
		}
	}

	return nil
}

func createLicenseIndexes(db *gorm.DB) error {
	log.Println("Creating license system indexes...")

	// License indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_license_keys_key_active ON license_keys(key, is_active)",
		"CREATE INDEX IF NOT EXISTS idx_license_keys_expires_at ON license_keys(expires_at) WHERE expires_at IS NOT NULL",
		"CREATE INDEX IF NOT EXISTS idx_user_licenses_user_active ON user_licenses(user_id, is_active)",
		"CREATE INDEX IF NOT EXISTS idx_user_licenses_license_active ON user_licenses(license_key_id, is_active)",
		"CREATE INDEX IF NOT EXISTS idx_user_licenses_activated_by ON user_licenses(activated_by, is_active)",
		"CREATE INDEX IF NOT EXISTS idx_user_licenses_last_used ON user_licenses(last_used) WHERE last_used IS NOT NULL",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			log.Printf("Failed to create license index: %s, Error: %v", indexSQL, err)
		}
	}

	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate is a public function that can be called from migration CLI
func AutoMigrate(db *gorm.DB) error {
	return autoMigrate(db)
}

func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}