package database

import (
	"backend/internal/services/database/service"
	"backend/internal/services/database/store"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"net"
)

type Config struct {
	Username string `required:"true"`
	Password string `required:"true"`
	Schema   string `required:"true"`
	Host     string `required:"true"`
	Port     int    `required:"true"`
	GrpcPort int    `required:"true"`
}

func NewService() (*service.Service, error) {
	var config Config
	if err := envconfig.Process("database", &config); err != nil {
		return nil, err
	}

	connection, err := store.OpenConnection(config.Username, config.Password, config.Schema, config.Host, config.Port)
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.GrpcPort))
	if err != nil {
		return nil, err
	}

	s := &service.Service{
		Listener: lis,
		DB:       connection,
	}
	return s, nil
}
