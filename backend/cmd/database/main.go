package main

import (
	"backend/internal/services"
	"backend/internal/services/database"
	"log"
)

func main() {
	s := database.NewService()

	if err := services.StartService(s); err != nil {
		log.Panic(err)
	}
}
