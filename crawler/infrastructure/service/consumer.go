package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
)

type daprMessageSubscriber struct {
	messages <-chan model.CrawlWebsite
}

func (f *daprMessageSubscriber) Consume(c context.Context) <-chan model.CrawlWebsite {
	return f.messages
}

func NewMessageConsumer(messages <-chan model.CrawlWebsite) service.MessageConsumer {
	return &daprMessageSubscriber{messages: messages}
}
