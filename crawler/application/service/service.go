package service

import (
	"context"
	"crawler/domain/model"
)

type MessageConsumer interface {
	Consume(c context.Context) chan *model.CrawlWebsite
}

type MessagePublisher interface {
	Publish(c context.Context, msg *model.CrawledWebsite) error
}

type WebsiteParser interface {
	Parse(url string) (*model.CrawledWebsite, error)
}
