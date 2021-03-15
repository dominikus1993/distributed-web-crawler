package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type ampqMessageConsumer struct {
	rabbitmq *amqp.Connection
	logger   *log.Logger
}

func subscribe(rabbitmq *amqp.Connection, logger *log.Logger, stream chan model.CrawlWebsite) {
	const exchange = "crawl-media"
	ch, err := rabbitmq.Channel()
	if err != nil {
		logger.WithError(err).Fatalln("Error when trying create channel")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"crawl-media", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		logger.WithError(err).Fatalln("Error when trying declare exchange")
	}

	q, err := ch.QueueDeclare(
		exchange, // name
		true,     // durable
		false,    // delete when usused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		logger.WithError(err).Fatalln("Error when trying create queue")
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		"#",      // routing key
		exchange, // exchange
		false,
		nil,
	)

	if err != nil {
		logger.WithError(err).Fatalln("Error when trying bind queue")
	}

	msgs, err := ch.Consume(
		q.Name,    // queue
		"crawler", // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	for msg := range msgs {
		var res model.CrawlWebsite
		err := json.Unmarshal(msg.Body, &res)
		if err != nil {
			log.Println(err)
		} else {
			stream <- res
		}
	}

	close(stream)
}

func (f *ampqMessageConsumer) Consume(c context.Context) chan model.CrawlWebsite {
	stream := make(chan model.CrawlWebsite)

	go subscribe(f.rabbitmq, f.logger, stream)

	return stream
}

func NewMessageConsumer(rabbitmq *amqp.Connection, logger *log.Logger) service.MessageConsumer {
	return &ampqMessageConsumer{rabbitmq: rabbitmq, logger: logger}
}
