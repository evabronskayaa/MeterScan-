package service

import (
	"backend/internal/mail"
	"backend/internal/proto"
	"context"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
	"log"
)

type Service struct {
	Mail     *mail.Client
	RabbitMQ *amqp.Connection
	Cron     *gocron.Scheduler
	DB       proto.DatabaseServiceClient
}

func (s *Service) Start() error {
	if _, err := s.Cron.Do(s.NotifyUsers); err != nil {
		return err
	}
	s.Cron.StartAsync()

	defer s.RabbitMQ.Close()

	channelRabbitMQ, err := s.RabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	q, err := channelRabbitMQ.QueueDeclare("mail", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	messages, err := channelRabbitMQ.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	go func() {
		for message := range messages {
			var m mail.Message
			if err := json.Unmarshal(message.Body, &m); err != nil {
				log.Println(" > Error json: ", err)
				continue
			}

			switch m.Type {
			case mail.Html:
				if err := s.Mail.SendHtmlMessage(m.Subject, m.File, m.Data, m.To...); err != nil {
					log.Println(" > Error: ", err)
					continue
				}
			case mail.Plain:
				if err := s.Mail.SendPlainMessage(m.Subject, m.Message, m.To...); err != nil {
					log.Println(" > Error: ", err)
					continue
				}
			}
		}
	}()
	return nil
}

func (s *Service) Shutdown(_ context.Context) {
	s.Cron.Stop()
}
