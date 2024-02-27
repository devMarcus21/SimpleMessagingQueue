package handlerfunctions

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataContracts "github.com/devMarcus21/SimpleMessagingQueue/src/api/datacontracts"
	errorResponses "github.com/devMarcus21/SimpleMessagingQueue/src/api/errorresponses"
	asyncQueueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/asyncqueue"
	queueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"
	logging "github.com/devMarcus21/SimpleMessagingQueue/src/logging"

	"github.com/google/uuid"
)

func buildQueueMessageFromQueueMessageRequest(request dataContracts.QueueMessageRequest, epochCreatedTimestamp int64) queueUtils.QueueMessage {
	newMessageId := uuid.New()
	return dataContracts.ConvertQueueMessageRequestToQueueMessage(request, newMessageId, epochCreatedTimestamp)
}

func PushMessageOntoQueueHandler(requestContext HandlerRequestContext, asyncQueue asyncQueueUtils.AsyncQueueWrapper) {
	var queueMessageRequest dataContracts.QueueMessageRequest

	err := json.NewDecoder(requestContext.GetHttpBody()).Decode(&queueMessageRequest)

	if err != nil {
		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleResponse(buildErrorResponse(requestContext.RequestStartTime(), errorResponses.JsonUnmarshalError.Message(), map[string]any{}))
		return
	}

	queueMessage := buildQueueMessageFromQueueMessageRequest(queueMessageRequest, requestContext.RequestStartTime())
	requestContext.Logger().Info(logging.APIPush_MessagePushedToQueueService.Message(), logging.LogIota, logging.APIPush_MessagePushedToQueueService.String(), "NewMessageId", queueMessage.MessageId)

	asyncQueue.Offer(queueMessage)

	requestContext.HandleResponse(
		buildSuccessfulResponse(
			requestContext.RequestStartTime(),
			SuccessfulResponseMessage,
			map[string]any{"MessageId": queueMessage.MessageId.String()}))
}

func PopMessageFromQueueHandler(requestContext HandlerRequestContext, asyncQueue asyncQueueUtils.AsyncQueueWrapper) {
	queueMessage, valueInQueue := asyncQueue.Poll()

	if !valueInQueue {
		requestContext.Logger().Info(logging.APIPop_QueueIsEmptyNoMessagePulled.Message(), logging.LogIota, logging.APIPop_QueueIsEmptyNoMessagePulled.String())

		requestContext.HandleResponse(buildSuccessfulResponse(requestContext.RequestStartTime(), "Queue is empty", map[string]any{}))
		return
	}

	requestContext.Logger().Info(logging.APIPop_MessagePulledFromQueueService.Message(), logging.LogIota, logging.APIPop_MessagePulledFromQueueService.String(), "PulledMessageId", queueMessage.MessageId)

	isBatched, batchIndex := queueMessage.IsBatchedMessage()
	requestContext.Logger().Info("Batched message property data", "IsBatched", isBatched, "BatchedIndex", batchIndex)

	requestContext.HandleResponse(
		buildSuccessfulResponse(
			requestContext.RequestStartTime(),
			SuccessfulResponseMessage,
			map[string]any{"QueueMessage": queueMessage}))
}

func BatchPushQueueMessagesOntoQueueHandler(requestContext HandlerRequestContext, asyncQueue asyncQueueUtils.AsyncQueueWrapper) {
	var batchQueueMessageRequest dataContracts.BatchQueueMessageRequest

	err := json.NewDecoder(requestContext.GetHttpBody()).Decode(&batchQueueMessageRequest)
	if err != nil {
		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleResponse(buildErrorResponse(requestContext.epochRequestStartTime, errorResponses.JsonUnmarshalError.Message(), map[string]any{}))
		return
	}

	batchSize := len(batchQueueMessageRequest.Messages)
	batchSizeMessage := fmt.Sprintf(logging.APIPushBatch_BatchSize.Message(), batchSize)
	requestContext.Logger().Info(batchSizeMessage, logging.LogIota, logging.APIPushBatch_BatchSize.String())

	if batchSize == 0 {
		requestContext.Logger().Error(errorResponses.GivenEmptyBatchError.String(), errorResponses.ApiErrorIota, errorResponses.GivenEmptyBatchError)

		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleResponse(buildErrorResponse(requestContext.epochRequestStartTime, errorResponses.GivenEmptyBatchError.Message(), map[string]any{}))
		return
	}

	maxBatchSize := requestContext.Configuration().Batching.MaxBatchPushSize
	if batchSize > maxBatchSize {
		batchToBigErrorMessage := fmt.Sprintf(errorResponses.BatchSizeBiggerThanMaxBatchSizeError.Message(), batchSize, maxBatchSize)

		requestContext.Logger().Error(batchToBigErrorMessage, errorResponses.ApiErrorIota, errorResponses.BatchSizeBiggerThanMaxBatchSizeError)

		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleResponse(buildErrorResponse(requestContext.epochRequestStartTime, batchToBigErrorMessage, map[string]any{}))
		return
	}

	processedMessageIds := []uuid.UUID{}

	for batchIndex, request := range batchQueueMessageRequest.Messages {
		queueMessage := buildQueueMessageFromQueueMessageRequest(request, requestContext.epochRequestStartTime)
		queueMessage.MakeBatchedMessage(batchIndex)

		asyncQueue.Offer(queueMessage)

		processedMessageIds = append(processedMessageIds, queueMessage.MessageId)
	}

	requestContext.HandleResponse(
		buildSuccessfulResponse(
			requestContext.RequestStartTime(),
			fmt.Sprintf(SuccessfullyBatchMessage, batchSize),
			map[string]any{"MessageIds": processedMessageIds}))
}
