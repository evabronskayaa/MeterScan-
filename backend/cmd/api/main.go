package main

import (
	"backend/internal/services"
	"backend/internal/services/api"
)

func main() {
	apiService := api.NewService()

	if err := services.StartService(apiService); err != nil {
		return
	}
}
