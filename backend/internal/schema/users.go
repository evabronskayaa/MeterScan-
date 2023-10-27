package schema

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Email    string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password []byte `json:"-"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Predictions []Prediction `json:"-"`
}

type Prediction struct {
	ID            uint `gorm:"primarykey" json:"id"`
	UserID        uint
	MeterReadings string  `gorm:"type:varchar(20);not null" json:"meter_readings"`
	Metric        float32 `json:"metric"`
}
