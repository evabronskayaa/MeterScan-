package main

import (
	"backend/internal/config"
	"backend/internal/services"
	"backend/internal/services/api"
	"log"
)

func main() {
	var apiConfig config.ApiConfig
	if err := apiConfig.Load(); err != nil {
		log.Panic("Config api load: ", err)
		return
	}

	apiService, err := api.NewService(apiConfig)
	if err != nil {
		log.Panic("Get api service: ", err)
		return
	}

	if err := services.StartService(apiService); err != nil {
		return
	}
}
