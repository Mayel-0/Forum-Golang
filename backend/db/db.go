package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func ConnectDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL manquante")
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("erreur ouverture GORM: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping échoué: %w", err)
	}

	log.Println("✅ DB Supabase connectée avec GORM!")

	return nil
}
