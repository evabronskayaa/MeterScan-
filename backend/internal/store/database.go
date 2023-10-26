package store

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenConnection(username, password, host, schema string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, schema)
	if conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		if err := migrate(conn); err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate()
}

func ShutdownConnection(db *gorm.DB) error {
	if sqlDB, err := db.DB(); err != nil {
		return err
	} else {
		return sqlDB.Close()
	}
}
