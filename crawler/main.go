package main

import (
	"context"
	"crawler/application/usecase"
	"crawler/infrastructure/env"
	"crawler/infrastructure/service"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func createLogger() *log.Logger {
	logger := log.New()
	logger.Formatter = &log.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}

	return logger
}

func main() {
	logger := createLogger()
	log.Infoln("Start Service")
	connection := env.GetEnvOrDefault("RABBITMQ_CONNECTION", "amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(connection)
	if err != nil {
		log.WithError(err).Fatalln("Failed to connect to RabbitMQ")
	}
	defer conn.Close()
	parser := service.NewWebsiteParser()
	publisher := service.NewMessagePublisher(conn)
	consumer := service.NewMessageConsumer(conn, logger)
	usecase := usecase.NewCrawlerUseCase(parser, publisher, consumer)
	ctx := context.TODO()
	usecase.StartCrawling(ctx)
	log.Infoln("Stop Service")
}
