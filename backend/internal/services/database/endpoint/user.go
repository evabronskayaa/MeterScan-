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

func (s GRPCServer) GetUser(_ context.Context, request *proto.UserRequest) (*proto.UserResponse, error) {
	if request.Id == nil && request.Email == nil {
		return nil, errors.ErrIncorrectRequest
	}

	query := "id = ?"
	var arg interface{} = request.Id
	if request.Email != nil {
		query = "email = ?"
		arg = request.Email
	}
	if user, err := getUser(s.DB, query, arg); err != nil {
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
