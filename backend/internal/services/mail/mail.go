package mail

import (
	"backend/internal/mail"
	"backend/internal/proto"
	"backend/internal/rabbit"
	"backend/internal/services/mail/service"
	"github.com/go-co-op/gocron"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Config struct {
	AmqpServerURL   string `required:"true"`
	DatabaseService string `required:"true"`

	Server   string `required:"true"`
	Port     int    `required:"true"`
	Login    string `required:"true"`
	Password string `required:"true"`
}

func NewService() *service.Service {
	var config Config
	if err := envconfig.Process("mail", &config); err != nil {
		panic(err)
	}

	mailClient, err := mail.NewMailClient(config.Server, config.Port, config.Login, config.Password)
	if err != nil {
		panic(err)
	} else if err := mailClient.Start(); err != nil {
		panic(err)
	}

	connectRabbitMQ, err := rabbit.ConnectRabbitMQ(config.AmqpServerURL)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	nextSchedule := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	cron := gocron.NewScheduler(now.Location()).
		Every(1).
		Hour().
		StartAt(nextSchedule)

	grpcConn, err := grpc.Dial(config.DatabaseService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	databaseClient := proto.NewDatabaseServiceClient(grpcConn)

	return &service.Service{
		Mail:     mailClient,
		RabbitMQ: connectRabbitMQ,
		Cron:     cron,
		DB:       databaseClient,
	}
}
