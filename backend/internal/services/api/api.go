package api

import (
	"backend/internal/config"
	"backend/internal/proto"
	"backend/internal/services/api/router"
	"backend/internal/services/api/service"
	"backend/internal/util"
	"github.com/go-co-op/gocron"
	"github.com/morkid/paginate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func NewService(config config.ApiConfig) (*service.Service, error) {
	conn, err := grpc.Dial(config.GRPCServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	grpcConn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	databaseClient := proto.NewDatabaseServiceClient(grpcConn)

	s := &service.Service{
		ReCaptcha:              reCaptcha,
		JWTSecret:              []byte(config.JWTSecret),
		ImageProcessingService: client,
		Pagination:             paginate.New(),
		Cron:                   cron,
		DatabaseService:        databaseClient,
	}
	s.Router = router.ConfigureRouter(s, config.Port)
	return s, nil
}
