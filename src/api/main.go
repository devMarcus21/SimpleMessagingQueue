package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dataContracts "github.com/devMarcus21/SimpleMessagingQueue/src/api/datacontracts"
	errorResponses "github.com/devMarcus21/SimpleMessagingQueue/src/api/errorresponses"
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

	SuccessfullyBatchMessage string = "Successfully added batch of size: %d"
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

func BuildSuccessfulPopResponse(time int64, responseMessage string, payload map[string]any) HttpServiceResponse {
	return HttpServiceResponse{
		Response:           SuccessfulResponse,
		RequestStartedTime: time,
		Message:            responseMessage,
		Payload:            payload,
	}
}

func buildQueueMessageFromQueueMessageRequest(request dataContracts.QueueMessageRequest) queueUtils.QueueMessage {
	newMessageId := uuid.New()
	epochTimeNow := time.Now().Unix()
	return dataContracts.ConvertQueueMessageRequestToQueueMessage(request, newMessageId, epochTimeNow)
}

func BuildHttpPushOntoQueueHandler(loggerBuilder logging.LoggerBuilder, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
		requestId := uuid.New()
		logger := loggerBuilder().With("RequestId", requestId)

		epochTimeStarted := time.Now().Unix()

		var queueMessageRequest dataContracts.QueueMessageRequest

		writer.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(reader.Body).Decode(&queueMessageRequest)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest) // TODO fix error responses to be json instead of just text
			return
		}

		queueMessage := buildQueueMessageFromQueueMessageRequest(queueMessageRequest)

		logger.Info(logging.APIPush_MessagePushedToQueueService.Message(), logging.LogIota, logging.APIPush_MessagePushedToQueueService.String(), "NewMessageId", queueMessage.MessageId)

		asyncQueue.Offer(queueMessage)

		json.NewEncoder(writer).Encode(BuildSuccessfulPushResponse(queueMessage.MessageId, epochTimeStarted))
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

		json.NewEncoder(writer).Encode(BuildSuccessfulPopResponse(epochTimeNow, SuccessfulResponseMessage, map[string]any{"QueueMessage": queueMessage}))
	}
}

func BuildHttpBatchPushOntoQueueHandler(loggerBuilder logging.LoggerBuilder, asyncQueue asyncQueueUtils.AsyncQueueWrapper, maxBatchSize int) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
		requestId := uuid.New()
		logger := loggerBuilder().With("RequestId", requestId)

		epochTimeStarted := time.Now().Unix()

		var batchQueueMessageRequest dataContracts.BatchQueueMessageRequest

		writer.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(reader.Body).Decode(&batchQueueMessageRequest)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest) // TODO fix error responses to be json instead of just text
			return
		}

		batchSize := len(batchQueueMessageRequest.Messages)
		batchSizeMessage := fmt.Sprintf(logging.APIPushBatch_BatchSize.Message(), batchSize)
		logger.Info(batchSizeMessage, logging.LogIota, logging.APIPushBatch_BatchSize.String())

		if batchSize == 0 {
			logger.Error(errorResponses.GivenEmptyBatchError.String(), errorResponses.ApiErrorIota, errorResponses.GivenEmptyBatchError)
			http.Error(writer, errorResponses.GivenEmptyBatchError.String(), http.StatusBadRequest) // TODO fix error responses to be json instead of just text
			return
		}
		if batchSize > maxBatchSize {
			logger.Error(fmt.Sprintf(errorResponses.BatchSizeBiggerThanMaxBatchSizeError.Message(), batchSize, maxBatchSize),
				errorResponses.ApiErrorIota, errorResponses.BatchSizeBiggerThanMaxBatchSizeError)
			http.Error(writer, errorResponses.BatchSizeBiggerThanMaxBatchSizeError.String(), http.StatusBadRequest) // TODO fix error responses to be json instead of just text
			return
		}

		processedMessageIds := []uuid.UUID{}
		for _, request := range batchQueueMessageRequest.Messages {
			queueMessage := buildQueueMessageFromQueueMessageRequest(request)
			asyncQueue.Offer(queueMessage)

			processedMessageIds = append(processedMessageIds, queueMessage.MessageId)
		}

		json.NewEncoder(writer).Encode(BuildSuccessfulPopResponse(epochTimeStarted, fmt.Sprintf(SuccessfullyBatchMessage, batchSize), map[string]any{"MessageIds": processedMessageIds}))
	}
}

func BuildHttpBatchPopFromQueueHandler(loggerBuilder logging.LoggerBuilder, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
	// TODO implement batch pop handler
	return func(writer http.ResponseWriter, reader *http.Request) {}
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

	http.HandleFunc("/push", BuildHttpPushOntoQueueHandler(loggerBuilder, asyncQueue))
	http.HandleFunc("/pop", BuildHttpPopFromQueueHandler(loggerBuilder, asyncQueue))

	http.HandleFunc("/push/batch", BuildHttpBatchPushOntoQueueHandler(loggerBuilder, asyncQueue, config.Batching.MaxBatchPushSize))
	http.HandleFunc("/pop/batch", BuildHttpBatchPopFromQueueHandler(loggerBuilder, asyncQueue))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
