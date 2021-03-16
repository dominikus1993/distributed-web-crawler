package main

import (
	log "github.com/sirupsen/logrus"
)

func createLogger() *log.Logger {
	logger := log.New()
	logger.Formatter = &log.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}
	return logger
}

func main() {
	logger := createLogger()
	logger.Infoln("Start Service")

}
