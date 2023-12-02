package api

import (
	"backend/internal/proto"
	"backend/internal/rabbit"
	"backend/internal/services/api/router"
	"backend/internal/services/api/service"
	"backend/internal/util"
	"github.com/go-co-op/gocron"
	"github.com/kelseyhightower/envconfig"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/morkid/paginate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Config struct {
	Port                 int    `required:"true"`
	JWTSecret            string `required:"true"`
	ReCaptchaSecret      string `required:"true"`
	GRPCServer           string `required:"true"`
	DatabaseService      string `required:"true"`
	RabbitMQ             string `required:"true"`
	MinioEndpoint        string `required:"true"`
	MinioAccessKey       string `required:"true"`
	MinioSecretAccessKey string `required:"true"`
}

func NewService() *service.Service {
	var config Config
	envconfig.MustProcess("api", &config)

	conn, err := grpc.Dial(config.GRPCServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := proto.NewImageProcessingServiceClient(conn)
	reCaptcha := util.ReCaptcha{
		Secret:  config.ReCaptchaSecret,
		Timeout: time.Second * 5,
	}

	rabbitmq, err := rabbit.ConnectRabbitMQ(config.RabbitMQ)
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

	s3Client, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretAccessKey, ""),
		Secure: true,
	})

	if err != nil {
		panic(err)
	}

	s := &service.Service{
		ReCaptcha:              reCaptcha,
		JWTSecret:              []byte(config.JWTSecret),
		ImageProcessingService: client,
		Pagination:             paginate.New(),
		Cron:                   cron,
		DatabaseService:        databaseClient,
		RabbitMQ:               rabbitmq,
		S3Client:               s3Client,
	}
	s.Router = router.ConfigureRouter(s, config.Port)
	return s
}
