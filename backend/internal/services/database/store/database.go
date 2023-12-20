package store

import (
	"backend/internal/services/database/schema"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func OpenConnection(username, password, schema, host string, port int) (*gorm.DB, error) {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
	})

	dsn := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", username, password, host, port, schema)
	if conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		return nil, err
	} else {
		if err := migrate(conn); err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(schema.User{}, schema.Prediction{}, schema.UserSetting{}, schema.PredictionInfo{})
}

func ShutdownConnection(db *gorm.DB) error {
	if sqlDB, err := db.DB(); err != nil {
		return err
	} else {
		return sqlDB.Close()
	}
}
