package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "github.com/devMarcus21/SimpleMessagingQueue/src/api/handlerfunctions"
	asyncQueueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/asyncqueue"
	configuration "github.com/devMarcus21/SimpleMessagingQueue/src/configuration"
	queueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"
	logging "github.com/devMarcus21/SimpleMessagingQueue/src/logging"
)

func main() {
	// Load service configuration
	config, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatal("Failed to read environment configuration file: ", err)
	}

	serviceLogger := logging.BuildTextLogger()
	serviceLogger.Info(
		"Starting queue service",
		"RunningOnPort", config.Port,
		"IsDevEnvironmentRunning", config.IsDevEnvironment,
		"LoggerType", config.Logging.LoggerType,
		"MaxBatchPushSize", config.Batching.MaxBatchPushSize,
		"MaxBatchReadSize", config.Batching.MaxBatchReadSize)

	var loggerBuilder logging.LoggerBuilder

	switch loggerType := config.Logging.LoggerType; loggerType {
	case "text":
		loggerBuilder = logging.BuildTextLogger
	case "json":
		loggerBuilder = logging.BuildJsonLogger
	default:
		serviceLogger.Warn("No or invalid logger type given in configuration (Defaulting to text logger)")
		loggerBuilder = logging.BuildTextLogger
	}

	queue := queueUtils.NewLinkedList()

	asyncQueue := asyncQueueUtils.NewAsyncQueue(queue)

	http.HandleFunc("/push", handlers.BuildHttpHandlerFunc(
		handlers.PushMessageOntoQueueHandler,
		loggerBuilder,
		asyncQueue,
		config,
	))

	http.HandleFunc("/pop", handlers.BuildHttpHandlerFunc(
		handlers.PopMessageFromQueueHandler,
		loggerBuilder,
		asyncQueue,
		config,
	))

	http.HandleFunc("/push/batch", handlers.BuildHttpHandlerFunc(
		handlers.BatchPushQueueMessagesOntoQueueHandler,
		loggerBuilder,
		asyncQueue,
		config,
	))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
