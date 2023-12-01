package service

import (
	"backend/internal/proto"
	"backend/internal/util"
	"context"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/morkid/paginate"
	"log"
	"net/http"
)

type Service struct {
	Router                 *http.Server
	ReCaptcha              util.ReCaptcha
	JWTSecret              []byte
	ImageProcessingService proto.ImageProcessingServiceClient
	Pagination             *paginate.Pagination
	Cron                   *gocron.Scheduler
	DatabaseService        proto.DatabaseServiceClient
}

func (s *Service) Start() error {
	go func() {
		if err := s.Router.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	if _, err := s.Cron.Do(s.NotifyUsers); err != nil {
		return err
	}
	s.Cron.StartAsync()
	return nil
}

func (s *Service) Shutdown(ctx context.Context) {
	if err := s.Router.Shutdown(ctx); err != nil {
		log.Printf("Web server shutdown with err: %v", err)
	}
	s.Cron.Stop()
}
