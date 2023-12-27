package endpoint

import (
	"backend/internal/errors"
	"backend/internal/proto"
	"backend/internal/services/database/schema"
	"backend/internal/util"
	"context"
	errors2 "errors"
	"gorm.io/gorm"
	"time"
)

type GRPCServer struct {
	proto.UnimplementedDatabaseServiceServer

	DB *gorm.DB
}

func getUser(db *gorm.DB, query string, args ...interface{}) (*schema.User, error) {
	var user *schema.User
	if db.Where(query, args).First(&user).Error != nil {
		return nil, errors.ErrNotFoundUser
	}
	return user, nil
}

func getUserByRequest(db *gorm.DB, request *proto.UserRequest) (*schema.User, error) {
	if request.Id == nil && request.Email == nil {
		return nil, errors.ErrIncorrectRequest
	}

	query := "id = ?"
	var arg interface{} = request.Id
	if request.Email != nil {
		query = "email = ?"
		arg = request.Email
	}

	return getUser(db, query, arg)
}

func (s GRPCServer) GetUser(_ context.Context, request *proto.UserRequest) (*proto.UserResponse, error) {
	if user, err := getUserByRequest(s.DB, request); err != nil {
		return nil, err
	} else {
		return user.Proto(), nil
	}
}

func (s GRPCServer) CreateUser(_ context.Context, request *proto.UserRequest) (*proto.UserResponse, error) {
	if request.Email != nil && request.Password != nil {
		email := *request.Email
		password := *request.Password

		if _, err := getUser(s.DB, "email = ?", email); !errors2.Is(err, errors.ErrNotFoundUser) {
			return nil, errors.ErrDuplicateEmail
		}

		hashPassword, err := util.HashPassword(password)
		if err != nil {
			return nil, errors.ErrCreatePassword
		}
		user := &schema.User{
			Email:    email,
			Password: string(hashPassword),
		}
		if s.DB.Create(&user).Error != nil {
			return nil, errors.ErrSaveUser
		} else {
			return user.Proto(), nil
		}
	}
	return nil, errors.ErrIncorrectRequest
}

func (s GRPCServer) VerifyUser(_ context.Context, request *proto.UserRequest) (*proto.Empty, error) {
	if request.Id != nil {
		user, err := getUser(s.DB, "id = ?", request.Id)
		if err != nil {
			return nil, err
		}

		if user.VerifiedAt.Valid {
			return nil, errors.ErrAlreadyVerified
		}

		if s.DB.Model(&user).Update("verified_at", time.Now()).Error != nil {
			return nil, errors.ErrSaveUser
		}

		return &proto.Empty{}, nil
	}
	return nil, errors.ErrIncorrectRequest
}

func (s GRPCServer) GetEmailsForNotification(_ context.Context, request *proto.GetEmailsForNotificationRequest) (*proto.EmailsResponse, error) {
	var users []schema.User
	dayOfMonth := request.Day
	hour := request.Hour

	if err := s.DB.Model(&schema.User{}).
		Preload("Settings").
		Joins("inner join user_settings on user_settings.user_id = users.id").
		Where("users.verified_at IS NOT NULL").
		Where("user_settings.notification_enabled = 1").
		Where("user_settings.notification_day_of_month = ?", dayOfMonth).
		Where("user_settings.notification_hour = ?", hour).
		Find(&users).Error; err != nil {
		return nil, errors.ErrNotFoundUser
	}

	var emails []string

	for _, user := range users {
		emails = append(emails, user.Email)
	}

	return &proto.EmailsResponse{
		Emails: emails,
	}, nil
}

func (s GRPCServer) GetSettings(_ context.Context, request *proto.UserRequest) (*proto.Settings, error) {
	var user schema.User

	result := s.DB.Preload("Settings").First(&user, request.Id)
	if result.Error != nil {
		if errors2.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.ErrNotFoundUser
		}
		return nil, result.Error
	}

	return user.Settings.Proto(), nil
}

func (s GRPCServer) UpdateSettings(_ context.Context, request *proto.UpdateSettingsRequest) (*proto.Empty, error) {
	var user schema.User

	result := s.DB.Preload("Settings").First(&user, request.Id)
	if result.Error != nil {
		if errors2.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.ErrNotFoundUser
		}
		return nil, result.Error
	}

	if user.Settings == (schema.UserSetting{}) {
		user.Settings = schema.UserSetting{UserID: user.ID}
		if err := s.DB.Create(&user.Settings).Error; err != nil {
			return nil, err
		}
	}

	settings := request.Settings
	if settings.NotificationDayOfMonth != nil {
		user.Settings.Notification.DayOfMonth = *settings.NotificationDayOfMonth
	}
	if settings.NotificationHour != nil {
		user.Settings.Notification.Hour = *settings.NotificationHour
	}
	if settings.NotificationEnabled != nil {
		user.Settings.Notification.Enabled = *settings.NotificationEnabled
	}

	if err := s.DB.Save(&user.Settings).Error; err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}
