package main

import (
	apiservice "backend/internal/api-service"
	"backend/internal/config"
	"backend/internal/service"
	"log"
)

func main() {
	var apiConfig config.ApiConfig
	if err := apiConfig.Load(); err != nil {
		log.Panic("Config api load: ", err)
		return
	}

	apiService, err := apiservice.NewService(apiConfig)
	if err != nil {
		log.Panic("Get api service: ", err)
		return
	}

	if err := service.StartService(apiService); err != nil {
		return
	}
}
