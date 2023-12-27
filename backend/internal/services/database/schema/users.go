package schema

import (
	"backend/internal/proto"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uint64 `gorm:"primarykey" json:"id"`
	Email    string `gorm:"type:varchar(254);not null;unique;index" json:"email"`
	Password string `gorm:"size:60" json:"-"`

	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	VerifiedAt sql.NullTime   `gorm:"default:null" json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Settings UserSetting `json:"settings"`

	Predictions []Prediction `gorm:"foreignKey:UserID" json:"-"`
}

type UserSetting struct {
	ID           uint64              `gorm:"primarykey" json:"-"`
	UserID       uint64              `json:"-"`
	Notification NotificationSetting `gorm:"embedded;embeddedPrefix:notification_" json:"notification"`
}

type NotificationSetting struct {
	Enabled    bool   `json:"enabled"`
	DayOfMonth uint32 `gorm:"type:int;not null" json:"day_of_month"`
	Hour       uint32 `gorm:"type:int;not null" json:"hour"`
}

func (u *User) Proto() *proto.UserResponse {
	return &proto.UserResponse{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Verified: u.VerifiedAt.Valid,
	}
}

func (s *UserSetting) Proto() *proto.Settings {
	return &proto.Settings{
		NotificationEnabled:    &s.Notification.Enabled,
		NotificationDayOfMonth: &s.Notification.DayOfMonth,
		NotificationHour:       &s.Notification.Hour,
	}
}
