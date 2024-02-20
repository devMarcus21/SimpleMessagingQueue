package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dataContracts "github.com/devMarcus21/SimpleMessagingQueue/src/api/datacontracts"
	asyncQueueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/asyncqueue"
	configuration "github.com/devMarcus21/SimpleMessagingQueue/src/configuration"
	queueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"
	logging "github.com/devMarcus21/SimpleMessagingQueue/src/logging"

	"github.com/google/uuid"
)

const (
	SuccessfulResponse        string = "Success"
	SuccessfulResponseMessage string = "Success"
	ErrorResponse             string = "Error"
)

type HttpServiceResponse struct {
	Response           string
	RequestStartedTime int64
	Message            string
	Payload            map[string]any
}

func BuildSuccessfulPushResponse(id uuid.UUID, time int64) HttpServiceResponse {
	return HttpServiceResponse{
		Response:           SuccessfulResponse,
		RequestStartedTime: time,
		Message:            SuccessfulResponseMessage,
		Payload: map[string]any{
			"MessageId": id.String(),
		},
	}
}

func BuildQueueEmptyResponse(time int64) HttpServiceResponse {
	return HttpServiceResponse{
		Response:           SuccessfulResponse,
		RequestStartedTime: time,
		Message:            "Queue is empty",
		Payload:            map[string]any{},
	}
}

func BuildSuccessfulPopResponse(message queueUtils.QueueMessage, time int64) HttpServiceResponse {
	return HttpServiceResponse{
		Response:           SuccessfulResponse,
		RequestStartedTime: time,
		Message:            SuccessfulResponseMessage,
		Payload: map[string]any{
			"QueueMessage": message,
		},
	}
}

func BuildHttpPushOntoQueueHandler(loggerBuilder logging.LoggerBuilder, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
		requestId := uuid.New()
		logger := loggerBuilder().With("RequestId", requestId)

		var QueueMessageRequest dataContracts.QueueMessageRequest

		writer.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(reader.Body).Decode(&QueueMessageRequest)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		newMessageId := uuid.New()
		epochTimeNow := time.Now().Unix()
		queueMessage := dataContracts.ConvertQueueMessageRequestToQueueMessage(QueueMessageRequest, newMessageId, epochTimeNow)

		logger.Info(logging.APIPush_MessagePushedToQueueService.Message(), logging.LogIota, logging.APIPush_MessagePushedToQueueService.String(), "NewMessageId", newMessageId)

		asyncQueue.Offer(queueMessage)

		json.NewEncoder(writer).Encode(BuildSuccessfulPushResponse(newMessageId, epochTimeNow))
	}
}

func BuildHttpPopFromQueueHandler(loggerBuilder logging.LoggerBuilder, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		requestId := uuid.New()
		logger := loggerBuilder().With("RequestId", requestId)

		epochTimeNow := time.Now().Unix()
		queueMessage, valueInQueue := asyncQueue.Poll()

		if !valueInQueue {
			logger.Info(logging.APIPop_QueueIsEmptyNoMessagePulled.Message(), logging.LogIota, logging.APIPop_QueueIsEmptyNoMessagePulled.String())
			json.NewEncoder(writer).Encode(BuildQueueEmptyResponse(epochTimeNow))
			return
		}

		logger.Info(logging.APIPop_MessagePulledFromQueueService.Message(), logging.LogIota, logging.APIPop_MessagePulledFromQueueService.String(), "PulledMessageId", queueMessage.MessageId)

		json.NewEncoder(writer).Encode(BuildSuccessfulPopResponse(queueMessage, epochTimeNow))
	}
}

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
		"IsLoggingEnabled", config.Logging.IsEnabled)

	loggerBuilder := logging.BuildEmptyLogger

	if config.Logging.IsEnabled {
		serviceLogger.Info("Logging is enabled", "LoggerType", config.Logging.LoggerType)

		switch loggerType := config.Logging.LoggerType; loggerType {
		case "text":
			loggerBuilder = logging.BuildTextLogger
		case "json":
			loggerBuilder = logging.BuildJsonLogger
		default:
			serviceLogger.Warn("No or invalid logger type given in configuration (logging is now turned off)")
		}
	}

	queue := queueUtils.NewLinkedList()

	asyncQueue := asyncQueueUtils.NewAsyncQueue(queue)

	http.HandleFunc("/push", BuildHttpPushOntoQueueHandler(loggerBuilder, asyncQueue))
	http.HandleFunc("/pop", BuildHttpPopFromQueueHandler(loggerBuilder, asyncQueue))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
