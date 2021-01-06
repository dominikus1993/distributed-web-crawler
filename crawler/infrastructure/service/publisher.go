package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"encoding/json"

	"github.com/streadway/amqp"
)

type ampqPublisher struct {
	rabbitmq *amqp.Connection
}

func (f *ampqPublisher) Publish(c context.Context, msg *model.CrawledWebsite) error {
	const exchange = "media"
	ch, err := f.rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	res, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchange, // exchange
		"#",      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        res,
		})

	return err
}

func NewMessagePublisher(rabbitmq *amqp.Connection) service.MessagePublisher {
	return &ampqPublisher{rabbitmq: rabbitmq}
}
