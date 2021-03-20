package main

import (
	"context"
	"crawler/application/usecase"
	"crawler/domain/model"
	"crawler/infrastructure/env"
	"crawler/infrastructure/logging"
	"crawler/infrastructure/service"
	"encoding/json"
	"fmt"
	"net/http"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

type DaprSubscription struct {
	PubSubName string `json:"pubsubname"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

func createLogger() *log.Logger {
	logger := log.New()
	logger.Formatter = &log.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}
	return logger
}

func getDaprSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs := [1]DaprSubscription{{PubSubName: "pubsub", Topic: env.GetEnvOrDefault("DAPR_SUB_TOPIC", "crawl-website"), Route: "crawl"}}
	res, err := json.Marshal(subs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func subscribe(stream chan<- model.CrawlWebsite, logger *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infoln("Subscribe")
		var res model.CrawlWebsite
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			logger.WithContext(r.Context()).WithError(err).Errorln("Error when trying read model in subscrition")
		} else {
			stream <- res
		}
	}
}

func main() {
	logger := createLogger()
	logger.Infoln("Start App")
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	logger.Infoln("Dapr initizlized")
	stream := make(chan model.CrawlWebsite)
	topic := env.GetEnvOrDefault("DAPR_PUBLISH_TOPIC", "crawled-media")
	pubsubname := env.GetEnvOrDefault("DAPR_PUBSUB_NAME", "pubsub")
	parser := service.NewWebsiteParser()
	publisher := service.NewMessagePublisher(client, topic, pubsubname)
	consumer := service.NewMessageConsumer(stream)
	usecase := usecase.NewCrawlerUseCase(parser, publisher, consumer)
	ctx, cancelWorkers := context.WithCancel(context.Background())
	go usecase.StartCrawling(ctx)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logging.NewStructuredLogger(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))
	router.HandleFunc("/dapr/subscribe", getDaprSubscriptions)
	router.HandleFunc("/crawl", subscribe(stream, logger))

	logger.Infoln("Start Service")
	err = http.ListenAndServe(":5000", router)
	if err != nil {
		logger.WithError(err).Fatalln("App host terminated")
	}
	cancelWorkers()
	close(stream)
	logger.Infoln("Stop Service")
}
