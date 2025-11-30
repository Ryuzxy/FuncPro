package db

import (
    "fmt"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/ryuzxy/fucpro/internal/config"
    "github.com/ryuzxy/fucpro/pkg/komoditas"
    "github.com/ryuzxy/fucpro/pkg/price"
)

// InitDB initializes database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        cfg.DBHost,
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBName,
        cfg.DBPort,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("opening DB: %w", err)
    }
    
    // Auto migrate tables
    if err := db.AutoMigrate(
        &komoditas.Komoditas{},
        &price.Price{},
    ); err != nil {
        return nil, fmt.Errorf("migrating DB: %w", err)
    }
    
    return db, nil
}