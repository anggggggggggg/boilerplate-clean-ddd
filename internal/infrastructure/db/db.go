package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sqlDB *sql.DB

func ConnectDB(dbConfig *config.DBConfig) (*gorm.DB, error) {
	dsn := buildDSN(dbConfig)
	db, sql, err := connect(dsn, dbConfig.DBPooling)
	if err != nil {
		return nil, err
	}

	sqlDB = sql

	log.Printf("Success to connect database: %s", dbConfig.Name)

	return db, nil
}

func CloseDB() {
	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		} else {
			log.Println("Success to close database connection")
		}
	}
}

func buildDSN(dbConfig *config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
		dbConfig.SSLMode,
	)
}

func connect(dsn string, dbPoolingConfig *config.DBPooling) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	// Connection pool config
	sqlDB.SetMaxOpenConns(dbPoolingConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbPoolingConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbPoolingConfig.ConnMaxLifetime) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(dbPoolingConfig.ConnMaxIdleTime) * time.Minute)

	return db, sqlDB, nil
}
