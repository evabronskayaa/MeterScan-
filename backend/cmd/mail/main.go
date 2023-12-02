package main

import (
	"backend/internal/mail"
	"backend/internal/rabbit"
	"github.com/goccy/go-json"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	AmqpServerURL string `required:"true"`

	Server   string `required:"true"`
	Port     int    `required:"true"`
	Login    string `required:"true"`
	Password string `required:"true"`
}

func main() {
	var config Config
	if err := envconfig.Process("mail", &config); err != nil {
		panic(err)
	}

	mailClient, err := mail.NewMailClient(config.Server, config.Port, config.Login, config.Password)
	if err != nil {
		panic(err)
	} else if err := mailClient.Start(); err != nil {
		panic(err)
	}

	connectRabbitMQ, err := rabbit.ConnectRabbitMQ(config.AmqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
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

	forever := make(chan bool)

	go func() {
		for message := range messages {
			var m mail.Message
			if err := json.Unmarshal(message.Body, &m); err != nil {
				log.Println(" > Error json: ", err)
				continue
			}

			switch m.Type {
			case mail.Html:
				if err := mailClient.SendHtmlMessage(m.Subject, m.File, m.Data, m.To...); err != nil {
					log.Println(" > Error: ", err)
					continue
				}
			case mail.Plain:
				if err := mailClient.SendPlainMessage(m.Subject, m.Message, m.To...); err != nil {
					log.Println(" > Error: ", err)
					continue
				}
			}
		}
	}()

	<-forever
}
