package main

import (
	"backend/internal/services"
	"backend/internal/services/mail"
	"log"
)

func main() {
	s := mail.NewService()

	if err := services.StartService(s); err != nil {
		log.Panic(err)
	}
}
