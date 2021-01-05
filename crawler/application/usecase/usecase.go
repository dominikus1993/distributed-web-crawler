package usecase

import (
	"context"
	"crawler/application/service"
	"crawler/domain/model"
	"log"
)

type CrawlerUseCase interface {
	StartCrawling(c context.Context)
}

type crawlerUseCase struct {
	parser     service.WebsiteParser
	publisher  service.MessagePublisher
	subscriber service.MessageConsumer
}

func publish(context context.Context, parsedChannel chan model.CrawledWebsite, publisher service.MessagePublisher) {
	for message := range parsedChannel {
		err := publisher.Publish(context, message)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func parse(msg model.CrawlWebsite, parser service.WebsiteParser, parsedChannel chan model.CrawledWebsite) {
	res, err := parser.Parse(msg.Url)
	if err != nil {
		log.Fatalln(err)
	}
	parsedChannel <- res
}

func (crawler *crawlerUseCase) StartCrawling(c context.Context) {
	crawledWebsitesStream := make(chan model.CrawledWebsite, 10)

	go publish(c, crawledWebsitesStream, crawler.publisher)

	consumeChannel := crawler.subscriber.Consume(c)

	for message := range consumeChannel {
		go parse(message, crawler.parser, crawledWebsitesStream)
	}

	close(crawledWebsitesStream)
}
