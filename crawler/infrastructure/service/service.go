package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"encoding/json"
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
		stream <- *model.NewCrawlWebsite("https://www.rossmann.pl/")
		close(stream)
	}()

	return stream
}

type htmlParser struct {
}

func (f *htmlParser) Parse(url string) (*model.CrawledWebsite, error) {
	contents := []model.Content{}
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL)
	})

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		contents = append(contents, model.NewContent(e.Attr("src")))
	})
	c.Visit(url)

	return model.NewCrawledWebsite(url, &contents), nil
}

type consolePublisher struct {
}

func (f *consolePublisher) Publish(c context.Context, msg *model.CrawledWebsite) error {
	res, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Println(string(res))
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
