package schema

import (
	"backend/internal/mail"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
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
	ID           uint                `gorm:"primarykey" json:"-"`
	UserID       uint                `json:"-"`
	Notification NotificationSetting `gorm:"embedded;embeddedPrefix:notification_" json:"notification"`
}

type NotificationSetting struct {
	DayOfMonth int `gorm:"type:int;not null" json:"day_of_month"`
	Hour       int `gorm:"type:int;not null" json:"hour"`
}

func (u *User) SendMessageToMail(client *mail.Client, subject, file string, data any) error {
	if !u.VerifiedAt.Valid {
		return nil
	}
	return client.SendHtmlMessage(subject, file, data, u.Email)
}

func (u *User) RequestVerification(client *mail.Client, token string) error {
	return u.SendMessageToMail(client, "Подтвердите почту", "request_verification.gohtml", token)
}
