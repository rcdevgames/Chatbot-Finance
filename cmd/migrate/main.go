package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"chatbot/internal/config"
	"chatbot/internal/database"
	"chatbot/internal/model"

	"gorm.io/gorm"
)

func main() {
	var (
		action      = flag.String("action", "up", "Migration action: up, down, status, reset")
		version     = flag.String("version", "", "Target version for rollback (use with down action)")
		force       = flag.Bool("force", false, "Force action without confirmation")
		configPath  = flag.String("config", ".env", "Path to config file")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Execute migration action
	switch *action {
	case "up":
		if err := migrateUp(db.DB); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
	case "down":
		if err := migrateDown(db.DB, *version, *force); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
	case "status":
		showMigrationStatus(db.DB)
	case "reset":
		if err := resetDatabase(db.DB, *force); err != nil {
			log.Fatalf("Database reset failed: %v", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, status, reset")
		os.Exit(1)
	}
}

func migrateUp(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Auto migrate will handle the migration logic
	return database.AutoMigrate(db)
}

func migrateDown(db *gorm.DB, targetVersion string, force bool) error {
	if !force {
		fmt.Print("⚠️  This will rollback migrations. Are you sure? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Migration cancelled")
			return nil
		}
	}

	fmt.Printf("Rolling back to version: %s\n", targetVersion)

	// Get all applied migrations
	var migrations []model.MigrationSchema
	if err := db.Where("applied = ?", true).Order("version DESC").Find(&migrations).Error; err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	for _, migration := range migrations {
		if targetVersion != "" && migration.Version <= targetVersion {
			break
		}

		fmt.Printf("Rolling back migration: %s (%s)\n", migration.Version, migration.Name)

		// Mark migration as not applied
		if err := db.Model(&migration).Update("applied", false).Error; err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
		}

		fmt.Printf("Migration %s rolled back successfully\n", migration.Version)
	}

	fmt.Println("Migration rollback completed")
	return nil
}

func showMigrationStatus(db *gorm.DB) {
	fmt.Println("Migration Status:")
	fmt.Println("================")

	// Get all migrations
	var migrations []model.MigrationSchema
	if err := db.Order("version").Find(&migrations).Error; err != nil {
		fmt.Printf("Failed to get migration status: %v\n", err)
		return
	}

	if len(migrations) == 0 {
		fmt.Println("No migrations found")
		return
	}

	for _, migration := range migrations {
		status := "❌ Not Applied"
		if migration.Applied {
			status = "✅ Applied"
		}
		fmt.Printf("%s - %s (%s)\n", migration.Version, migration.Name, status)
	}
}

func resetDatabase(db *gorm.DB, force bool) error {
	if !force {
		fmt.Print("⚠️  This will DROP ALL TABLES and recreate them. Are you sure? (type 'yes' to confirm): ")
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			fmt.Println("Database reset cancelled")
			return nil
		}
	}

	fmt.Println("Dropping all tables...")

	// Get all table names
	var tables []string
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	// Drop all tables
	for _, table := range tables {
		if table == "migration_schemas" {
			continue // Skip migration table
		}
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	// Reset migration table
	if err := db.Exec("DELETE FROM migration_schemas").Error; err != nil {
		return fmt.Errorf("failed to reset migration schemas: %w", err)
	}

	fmt.Println("Running fresh migrations...")
	return database.AutoMigrate(db)
}