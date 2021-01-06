package service

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"log"

	"github.com/gocolly/colly/v2"
)

type fakeMessageConsumer struct {
}

func (f *fakeMessageConsumer) Consume(c context.Context) chan model.CrawlWebsite {
	stream := make(chan model.CrawlWebsite)

	go func() {
		stream <- model.NewCrawlWebsite("https://jbzd.com.pl/")
		stream <- model.NewCrawlWebsite("https://jbzd.com.pl/")
		stream <- model.NewCrawlWebsite("https://jbzd.com.pl/")
		stream <- model.NewCrawlWebsite("https://jbzd.com.pl/")
		stream <- model.NewCrawlWebsite("https://jbzd.com.pl/")
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

func NewMessageConsumer() service.MessageConsumer {
	return &fakeMessageConsumer{}
}

func NewWebsiteParser() service.WebsiteParser {
	return &htmlParser{}
}
