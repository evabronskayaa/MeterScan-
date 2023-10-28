package api_service

import (
	"backend/internal/api-service/service"
	"backend/internal/config"
	"backend/internal/proto"
	"backend/internal/store"
	"backend/internal/util"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Service service.Service

func NewService(config config.ApiConfig) (*Service, error) {
	dbConfig := config.Database
	connection, err := store.OpenConnection(dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Schema)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(config.GRPCServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewImageProcessingServiceClient(conn)

	s := &Service{
		Port: config.Port,
		DB:   connection,
		ReCaptcha: util.ReCaptcha{
			Secret:  config.ReCaptchaSecret,
			Timeout: time.Second * 5,
		},
		JWTSecret:              []byte(config.JWTSecret),
		ImageProcessingService: client,
	}
	s.Router = s.configureRouter()
	return s, nil
}

func (s *Service) Start() error {
	go func() {
		if err := s.Router.Run(fmt.Sprintf(":%v", s.Port)); err != nil {
			panic(err)
		}
	}()
	return nil
}

func (s *Service) Shutdown() error {
	return store.ShutdownConnection(s.DB)
}
