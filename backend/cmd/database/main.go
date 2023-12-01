package main

import (
	"backend/internal/services"
	"backend/internal/services/database"
	"log"
)

func main() {
	if s, err := database.NewService(); err != nil {
		log.Panic("Get database service: ", err)
	} else if err := services.StartService(s); err != nil {
		log.Panic(err)
	}
}
