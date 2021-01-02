package service

import "crawler/domain/model"

type MessageConsumer interface {
	Consume() chan model.CrawledWebsite
}
