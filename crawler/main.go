package main

import (
	"context"
	"crawler/application/usecase"
	"crawler/infrastructure/service"
)

func main() {
	parser := service.NewWebsiteParser()
	publisher := service.NewMessagePublisher()
	consumer := service.NewMessageConsumer()
	usecase := usecase.NewCrawlerUseCase(parser, publisher, consumer)
	ctx := context.TODO()
	usecase.StartCrawling(ctx)
	print("test")
}
