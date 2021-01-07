package main

import (
	"context"
	"crawler/application/usecase"
	"crawler/infrastructure/env"
	"crawler/infrastructure/service"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	log.Println("Start Service")
	connection := env.GetEnvOrDefault("RABBITMQ_CONNECTION", "amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(connection)
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()
	parser := service.NewWebsiteParser()
	publisher := service.NewMessagePublisher(conn)
	consumer := service.NewMessageConsumer(conn)
	usecase := usecase.NewCrawlerUseCase(parser, publisher, consumer)
	ctx := context.TODO()
	usecase.StartCrawling(ctx)
	log.Println("Stop Service")
}
