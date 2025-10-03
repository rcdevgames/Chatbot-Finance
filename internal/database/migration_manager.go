package database

import (
	"fmt"
	"log"

	"chatbot/internal/model"

	"gorm.io/gorm"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db *gorm.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// CreateMigration creates a new migration with the given version and models
func (m *MigrationManager) CreateMigration(version, name string, models ...interface{}) error {
	// Check if migration already exists
	var existing model.MigrationSchema
	if err := m.db.Where("version = ?", version).First(&existing).Error; err == nil {
		return fmt.Errorf("migration with version %s already exists", version)
	}

	// Create migration record
	migration := model.MigrationSchema{
		Version: version,
		Name:    name,
		Applied: false,
	}

	if err := m.db.Create(&migration).Error; err != nil {
		return fmt.Errorf("failed to create migration record: %w", err)
	}

	log.Printf("Created migration: %s - %s", version, name)
	return nil
}

// RunMigration runs a specific migration
func (m *MigrationManager) RunMigration(version string) error {
	var migration model.MigrationSchema
	if err := m.db.Where("version = ?", version).First(&migration).Error; err != nil {
		return fmt.Errorf("migration %s not found: %w", version, err)
	}

	if migration.Applied {
		return fmt.Errorf("migration %s is already applied", version)
	}

	log.Printf("Running migration: %s - %s", migration.Version, migration.Name)

	// Apply the migration
	if err := autoMigrate(m.db); err != nil {
		return fmt.Errorf("failed to run migration %s: %w", version, err)
	}

	// Mark as applied
	if err := m.db.Model(&migration).Update("applied", true).Error; err != nil {
		return fmt.Errorf("failed to mark migration %s as applied: %w", version, err)
	}

	log.Printf("Migration %s completed successfully", version)
	return nil
}

// GetPendingMigrations returns all pending migrations
func (m *MigrationManager) GetPendingMigrations() ([]model.MigrationSchema, error) {
	var migrations []model.MigrationSchema
	if err := m.db.Where("applied = ?", false).Order("version").Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending migrations: %w", err)
	}
	return migrations, nil
}

// GetAppliedMigrations returns all applied migrations
func (m *MigrationManager) GetAppliedMigrations() ([]model.MigrationSchema, error) {
	var migrations []model.MigrationSchema
	if err := m.db.Where("applied = ?", true).Order("version").Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	return migrations, nil
}

// RollbackToVersion rolls back to a specific version
func (m *MigrationManager) RollbackToVersion(targetVersion string) error {
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	for i := len(appliedMigrations) - 1; i >= 0; i-- {
		migration := appliedMigrations[i]
		if migration.Version <= targetVersion {
			break
		}

		log.Printf("Rolling back migration: %s - %s", migration.Version, migration.Name)

		// Mark migration as not applied
		if err := m.db.Model(&migration).Update("applied", false).Error; err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
		}

		log.Printf("Migration %s rolled back successfully", migration.Version)
	}

	return nil
}

// ValidateMigrations checks if all migrations are in a consistent state
func (m *MigrationManager) ValidateMigrations() error {
	// Check for duplicate versions
	var duplicates []struct {
		Version string
		Count   int
	}

	if err := m.db.Raw("SELECT version, COUNT(*) as count FROM migration_schemas GROUP BY version HAVING COUNT(*) > 1").Scan(&duplicates).Error; err != nil {
		return fmt.Errorf("failed to check for duplicate migrations: %w", err)
	}

	if len(duplicates) > 0 {
		return fmt.Errorf("found duplicate migration versions: %v", duplicates)
	}

	// Check if all required tables exist
	requiredTables := []string{"users", "categories", "transactions", "chat_history", "budgets", "insights"}
	for _, table := range requiredTables {
		var count int64
		if err := m.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", table).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if count == 0 {
			return fmt.Errorf("required table %s does not exist", table)
		}
	}

	log.Println("Migration validation passed")
	return nil
}

// GetMigrationHistory returns the complete migration history
func (m *MigrationManager) GetMigrationHistory() ([]model.MigrationSchema, error) {
	var migrations []model.MigrationSchema
	if err := m.db.Order("version").Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get migration history: %w", err)
	}
	return migrations, nil
}