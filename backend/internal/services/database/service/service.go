package service

import (
	"backend/internal/proto"
	"backend/internal/services/database/endpoint"
	"backend/internal/services/database/store"
	"context"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
)

type Service struct {
	Listener net.Listener
	DB       *gorm.DB
}

func (s *Service) Start() error {
	server := grpc.NewServer()
	proto.RegisterDatabaseServiceServer(server, endpoint.GRPCServer{DB: s.DB})
	log.Printf("server listening at %v", s.Listener.Addr())
	if err := server.Serve(s.Listener); err != nil {
		return err
	}
	return nil
}

func (s *Service) Shutdown(_ context.Context) {
	if err := store.ShutdownConnection(s.DB); err != nil {
		log.Printf("Database shutdown with err: %v", err)
	}
}
