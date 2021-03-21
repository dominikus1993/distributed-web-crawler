package service

import (
	"crawler/application/service"
	"crawler/domain/model"
	"log"

	"github.com/gocolly/colly/v2"
)

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

func NewWebsiteParser() service.WebsiteParser {
	return &htmlParser{}
}
