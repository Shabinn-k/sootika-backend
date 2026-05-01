package database

import (
	"fmt"
	"golang/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

var (
	pgDB   *gorm.DB
	pgOnce sync.Once
)

func SetupDatabase(cfg *config.Config) *gorm.DB {

	pgOnce.Do(func() {

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			cfg.DB.Host,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Name,
			cfg.DB.Port,
			cfg.DB.SSLMode,
			cfg.DB.TimeZone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("Failed to get DB instance:", err)
		}

		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

		if err := sqlDB.Ping(); err != nil {
			log.Fatal("Database not reachable:", err)
		}

		pgDB = db
		log.Println("Database connected successfully")
	})

	return pgDB
}
