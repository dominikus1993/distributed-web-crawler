package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

type fakeMessageConsumer struct {
}

func (f *fakeMessageConsumer) Consume(c context.Context) chan model.CrawlWebsite {
	stream := make(chan model.CrawlWebsite)

	go func() {
		stream <- *model.NewCrawlWebsite("https://jbzd.com.pl/")
		time.Sleep(2 * time.Second)
		stream <- *model.NewCrawlWebsite("https://httpbin.org/delay/1")
		close(stream)
	}()

	return stream
}

type htmlParser struct {
}

func (f *htmlParser) Parse(url string) (*model.CrawledWebsite, error) {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})
	c.Visit(url)
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
	return &htmlParser{}
}

func NewMessagePublisher() service.MessagePublisher {
	return &consolePublisher{}
}
