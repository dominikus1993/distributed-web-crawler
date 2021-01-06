package main

import (
	"context"
	"crawler/application/usecase"
	"crawler/infrastructure/service"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
	print("test")
}
