package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"encoding/json"

	dapr "github.com/dapr/go-sdk/client"
)

type daprPublisher struct {
	client      dapr.Client
	topic       string
	pubsubaname string
}

func (f *daprPublisher) Publish(c context.Context, msg *model.CrawledWebsite) error {
	res, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = f.client.PublishEvent(c, f.pubsubaname, f.topic, res)
	return err
}

func NewMessagePublisher(client dapr.Client, topic string, pubsubname string) service.MessagePublisher {
	return &daprPublisher{client: client, topic: topic, pubsubaname: pubsubname}
}
