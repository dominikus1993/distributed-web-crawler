package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"log"
)

type fakeMessageConsumer struct {
}

func (f *fakeMessageConsumer) Consume(c context.Context) chan *model.CrawlWebsite {
	stream := make(chan *model.CrawlWebsite)

	go func() {
		res := model.NewCrawlWebsite("https://jbzd.com.pl/")
		stream <- res
	}()

	return stream
}

type fakeParser struct {
}

func (f *fakeParser) Parse(url string) (*model.CrawledWebsite, error) {
	return model.NewCrawledWebsite(url), nil
}

type consolePublisher struct {
}

func (f *consolePublisher) Publish(c context.Context, msg *model.CrawledWebsite) error {
	log.Println(msg)
	return nil
}

func NewMessageConsumer() service.MessageConsumer {
	return &fakeMessageConsumer{}
}

func NewWebsiteParser() service.WebsiteParser {
	return &fakeParser{}
}

func NewMessagePublisher() service.MessagePublisher {
	return &consolePublisher{}
}
