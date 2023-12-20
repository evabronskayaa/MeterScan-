package service

import (
	"backend/internal/proto"
	"backend/internal/util"
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/morkid/paginate"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type Service struct {
	Router                 *http.Server
	ReCaptcha              util.ReCaptcha
	JWTSecret              []byte
	ImageProcessingService proto.ImageProcessingServiceClient
	Pagination             *paginate.Pagination
	DatabaseService        proto.DatabaseServiceClient
	RabbitMQ               *amqp.Connection
	S3Client               *minio.Client
}

func (s *Service) Start() error {
	go func() {
		if err := s.Router.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	return nil
}

func (s *Service) Shutdown(ctx context.Context) {
	if err := s.RabbitMQ.Close(); err != nil {
		log.Printf("RabbitMQ shutdown with err: %v", err)
	}
	if err := s.Router.Shutdown(ctx); err != nil {
		log.Printf("Web server shutdown with err: %v", err)
	}
}
