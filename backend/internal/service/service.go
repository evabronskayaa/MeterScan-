package service

import (
	"context"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

type Service interface {
	Start() error
	Shutdown(ctx context.Context)
}

func StartService(s Service) error {
	if err := s.Start(); err != nil {
		return err
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutdown service %v...", reflect.TypeOf(s))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Printf("Service %v exiting", reflect.TypeOf(s))
	return nil
}
