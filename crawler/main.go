package main

import (
	"context"
	"crawler/application/usecase"
	"crawler/domain/model"
	"crawler/infrastructure/env"
	"crawler/infrastructure/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
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
	logger.Infoln("Start Service")
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	stream := make(chan model.CrawlWebsite)
	log.Println("Start Service")
	topic := env.GetEnvOrDefault("DAPR_PUBLISH_TOPIC", "crawled")
	pubsubname := env.GetEnvOrDefault("DAPR_PUBSUB_NAME", "pubsub")
	parser := service.NewWebsiteParser()
	publisher := service.NewMessagePublisher(client, topic, pubsubname)
	consumer := service.NewMessageConsumer(stream)
	usecase := usecase.NewCrawlerUseCase(parser, publisher, consumer)
	ctx, cancelWorkers := context.WithCancel(context.Background())
	usecase.StartCrawling(ctx)
	router := mux.NewRouter()
	router.HandleFunc("/dapr/subscribe", getDaprSubscriptions)
	router.HandleFunc("/crawl", subscribe(stream, logger))
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.WithError(err).Fatalln("App host terminated")
	}
	cancelWorkers()
	close(stream)
	log.Println("Stop Service")
}
