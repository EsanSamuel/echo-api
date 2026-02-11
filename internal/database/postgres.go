package database

import (
	"github.com/echo/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.ConnectionUrl()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
