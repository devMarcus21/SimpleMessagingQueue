package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
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

func BuildHttpPushOntoQueueHandler(logger *slog.Logger, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
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

		logger.Info(logging.MessagePushedToQueueService.Message(), "NewMessageId", newMessageId)

		asyncQueue.Offer(queueMessage)

		json.NewEncoder(writer).Encode(BuildSuccessfulPushResponse(newMessageId, epochTimeNow))
	}
}

func BuildHttpPopFromQueueHandler(logger *slog.Logger, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		epochTimeNow := time.Now().Unix()
		queueMessage, valueInQueue := asyncQueue.Poll()

		if !valueInQueue {
			json.NewEncoder(writer).Encode(BuildQueueEmptyResponse(epochTimeNow))
			return
		}

		logger.Info(logging.MessagePulledFromQueueService.Message(), "MessageId", queueMessage.MessageId)

		json.NewEncoder(writer).Encode(BuildSuccessfulPopResponse(queueMessage, epochTimeNow))
	}
}

func main() {
	// Load service configuration
	config, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatal("Failed to read environment configuration file: ", err)
	}

	fmt.Println("Dev environment running: ", config.IsDevEnvironment)

	//logger := logging.BuildEmptyLogger()
	logger := logging.BuildTextLogger()

	queue := queueUtils.NewLinkedList()

	asyncQueue := asyncQueueUtils.NewAsyncQueue(queue)

	http.HandleFunc("/push", BuildHttpPushOntoQueueHandler(logger, asyncQueue))
	http.HandleFunc("/pop", BuildHttpPopFromQueueHandler(logger, asyncQueue))

	log.Fatal(http.ListenAndServe(":80", nil))
}
