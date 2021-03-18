package main

import (
	"crawler/domain/model"
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

func hello(d dapr.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := model.CrawledWebsite{Url: "tets"}
		res, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error")
		}
		err = d.PublishEvent(r.Context(), "pubsub", "test", res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error dapr", err)
		}
		fmt.Fprintln(w, "Welcome!")
	}
}

func getDaprSubscriptions(w http.ResponseWriter, r *http.Request) {
	log.Infoln("AAAA")

	subs := [1]DaprSubscription{DaprSubscription{PubSubName: "pubsub", Topic: "test", Route: "testsub"}}
	res, err := json.Marshal(subs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	log.Infoln("BBBBB")
}

func main() {
	logger := createLogger()
	logger.Infoln("Start Service")
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	router := mux.NewRouter()
	router.HandleFunc("/hello", hello(client))
	router.HandleFunc("/dapr/subscribe", getDaprSubscriptions)
	router.HandleFunc("/testsub", subscribe)
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
