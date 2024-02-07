package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	dataContracts "github.com/devMarcus21/SimpleMessagingQueue/src/api/datacontracts"
	asyncQueueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/asyncqueue"
	queueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"
	logging "github.com/devMarcus21/SimpleMessagingQueue/src/logging"

	"github.com/google/uuid"
)

type HttpServiceResponse map[string]any

func BuildSuccessfulPushResponse(id uuid.UUID, time int64) HttpServiceResponse {
	return HttpServiceResponse{
		"status":           "success",
		"messageId":        id.String(),
		"epochTimeStarted": time,
	}
}

func BuildQueueEmptyResponse(time int64) HttpServiceResponse {
	return HttpServiceResponse{
		"status":           "error",
		"message":          "Queue is empty",
		"epochTimeStarted": time,
	}
}

func BuildSuccessfulPopResponse(message queueUtils.QueueMessage, time int64) HttpServiceResponse {
	return HttpServiceResponse{
		"status":           "success",
		"epochTimeStarted": time,
		"message":          message,
	}
}

func BuildHttpPushOntoQueueHandler(logger logging.Logger, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
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

		logger.Log(newMessageId, logging.MessagePushedToQueueService, newMessageId)

		asyncQueue.Offer(queueMessage)

		json.NewEncoder(writer).Encode(BuildSuccessfulPushResponse(newMessageId, epochTimeNow))
	}
}

func BuildHttpPopFromQueueHandler(logger logging.Logger, asyncQueue asyncQueueUtils.AsyncQueueWrapper) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		epochTimeNow := time.Now().Unix()
		queueMessage, valueInQueue := asyncQueue.Poll()

		if !valueInQueue {
			json.NewEncoder(writer).Encode(BuildQueueEmptyResponse(epochTimeNow))
			return
		}

		logger.Log(queueMessage.MessageId, logging.MessagePulledFromQueueService, queueMessage.MessageId)

		json.NewEncoder(writer).Encode(BuildSuccessfulPopResponse(queueMessage, epochTimeNow))
	}
}

func main() {
	logger := logging.NewEmptyLogger()
	queue := queueUtils.NewLinkedList()

	//var asyncQueue asyncQueueUtils.AsyncQueueWrapper
	asyncQueue := asyncQueueUtils.NewAsyncQueue(queue)

	http.HandleFunc("/push", BuildHttpPushOntoQueueHandler(logger, asyncQueue))
	http.HandleFunc("/pop", BuildHttpPopFromQueueHandler(logger, asyncQueue))

	log.Fatal(http.ListenAndServe(":80", nil))
}
