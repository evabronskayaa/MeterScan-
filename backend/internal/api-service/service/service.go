package service

import (
	"backend/internal/mail"
	"backend/internal/proto"
	"backend/internal/store"
	"backend/internal/util"
	"context"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Service struct {
	Router                 *http.Server
	DB                     *gorm.DB
	ReCaptcha              util.ReCaptcha
	JWTSecret              []byte
	MailClient             *mail.Client
	ImageProcessingService proto.ImageProcessingServiceClient
	Pagination             *paginate.Pagination
	Cron                   *gocron.Scheduler
}

func (s *Service) Start() error {
	if err := s.MailClient.Start(); err != nil {
		return err
	}
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
	s.MailClient.Shutdown(ctx)
	if err := store.ShutdownConnection(s.DB); err != nil {
		log.Printf("Shutdown database with err: %v", err)
	}
}
