package api_service

import (
	"backend/internal/api-service/router"
	"backend/internal/api-service/service"
	"backend/internal/config"
	"backend/internal/mail"
	"backend/internal/proto"
	"backend/internal/store"
	"backend/internal/util"
	"context"
	"github.com/go-co-op/gocron"
	"github.com/morkid/paginate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func NewService(config config.ApiConfig) (*service.Service, error) {
	dbConfig := config.Database
	connection, err := store.OpenConnection(dbConfig.Username, dbConfig.Password, dbConfig.Schema, dbConfig.Host, dbConfig.Port)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, config.GRPCServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	mailConfig := config.Mail
	mailClient, err := mail.NewMailClient(mailConfig.Server, mailConfig.Port, mailConfig.Login, mailConfig.Password)
	if err != nil {
		return nil, err
	}

	client := proto.NewImageProcessingServiceClient(conn)
	reCaptcha := util.ReCaptcha{
		Secret:  config.ReCaptchaSecret,
		Timeout: time.Second * 5,
	}

	now := time.Now()
	nextSchedule := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	cron := gocron.NewScheduler(now.Location()).
		Every(1).
		Hour().
		StartAt(nextSchedule)

	s := &service.Service{
		DB:                     connection,
		ReCaptcha:              reCaptcha,
		JWTSecret:              []byte(config.JWTSecret),
		ImageProcessingService: client,
		MailClient:             mailClient,
		Pagination:             paginate.New(),
		Cron:                   cron,
	}
	s.Router = router.ConfigureRouter(s, config.Port)
	return s, nil
}
