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
		requestContext.Logger().Error(errorResponses.JsonUnmarshalError.String(), errorResponses.ApiErrorIota, errorResponses.JsonUnmarshalError.String())
		requestContext.Logger().Error(fmt.Sprintf(logging.JsonDecodeError.Message(), err.Error()), logging.LogEventIota, logging.JsonDecodeError.String())

		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleHttpResponse(buildErrorResponse(errorResponses.JsonUnmarshalError.Message(), map[string]any{}))
		return
	}

	queueMessage := buildQueueMessageFromQueueMessageRequest(queueMessageRequest, requestContext.RequestStartTime())
	requestContext.Logger().Info(logging.MessagePushedToQueueService.Message(), logging.LogEventIota, logging.MessagePushedToQueueService.String(), string(logging.NewMessageId), queueMessage.MessageId)

	asyncQueue.Offer(queueMessage)

	requestContext.HandleHttpResponse(buildSuccessfulResponse(SuccessfulResponseMessage, map[string]any{"MessageId": queueMessage.MessageId.String()}))
}

func PopMessageFromQueueHandler(requestContext HandlerRequestContext, asyncQueue asyncQueueUtils.AsyncQueueWrapper) {
	queueMessage, valueInQueue := asyncQueue.Poll()

	if !valueInQueue {
		requestContext.Logger().Info(logging.QueueIsEmptyNoMessagePulled.Message(), logging.LogEventIota, logging.QueueIsEmptyNoMessagePulled.String())
		requestContext.HandleHttpResponse(buildSuccessfulResponse(QueueIsEmptyMessage, map[string]any{}))
		return
	}

	requestContext.Logger().Info(logging.MessagePulledFromQueueService.Message(), logging.LogEventIota, logging.MessagePulledFromQueueService.String(), string(logging.PulledMessageId), queueMessage.MessageId)

	isBatched, batchIndex := queueMessage.IsBatchedMessage()
	requestContext.Logger().Info(logging.BatchMessageProperties.Message(), logging.LogEventIota, logging.BatchMessageProperties, string(logging.IsBatched), isBatched, string(logging.BatchedIndex), batchIndex)

	requestContext.HandleHttpResponse(buildSuccessfulResponse(SuccessfulResponseMessage, map[string]any{"QueueMessage": queueMessage}))
}

func BatchPushQueueMessagesOntoQueueHandler(requestContext HandlerRequestContext, asyncQueue asyncQueueUtils.AsyncQueueWrapper) {
	var batchQueueMessageRequest dataContracts.BatchQueueMessageRequest

	err := json.NewDecoder(requestContext.GetHttpBody()).Decode(&batchQueueMessageRequest)
	if err != nil {
		requestContext.Logger().Error(errorResponses.JsonUnmarshalError.String(), errorResponses.ApiErrorIota, errorResponses.JsonUnmarshalError.String())
		requestContext.Logger().Error(fmt.Sprintf(logging.JsonDecodeError.Message(), err.Error()), logging.LogEventIota, logging.JsonDecodeError)

		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleHttpResponse(buildErrorResponse(errorResponses.JsonUnmarshalError.Message(), map[string]any{}))
		return
	}

	batchSize := len(batchQueueMessageRequest.Messages)
	batchSizeMessage := fmt.Sprintf(logging.PushBatchBatchSize.Message(), batchSize)
	requestContext.Logger().Info(batchSizeMessage, logging.LogEventIota, logging.PushBatchBatchSize.String())

	if batchSize == 0 {
		requestContext.Logger().Error(errorResponses.GivenEmptyBatchError.String(), errorResponses.ApiErrorIota, errorResponses.GivenEmptyBatchError.String())

		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleHttpResponse(buildErrorResponse(errorResponses.GivenEmptyBatchError.Message(), map[string]any{}))
		return
	}

	maxBatchSize := requestContext.Configuration().Batching.MaxBatchPushSize
	if batchSize > maxBatchSize {
		batchToBigErrorMessage := fmt.Sprintf(errorResponses.BatchSizeBiggerThanMaxBatchSizeError.Message(), batchSize, maxBatchSize)

		requestContext.Logger().Error(batchToBigErrorMessage, errorResponses.ApiErrorIota, errorResponses.BatchSizeBiggerThanMaxBatchSizeError.String())

		requestContext.AddHttpStatusCode(http.StatusBadRequest)
		requestContext.HandleHttpResponse(buildErrorResponse(batchToBigErrorMessage, map[string]any{}))
		return
	}

	processedMessageIds := []uuid.UUID{}

	for batchIndex, request := range batchQueueMessageRequest.Messages {
		queueMessage := buildQueueMessageFromQueueMessageRequest(request, requestContext.epochRequestStartTime)
		queueMessage.MakeBatchedMessage(batchIndex)

		asyncQueue.Offer(queueMessage)

		processedMessageIds = append(processedMessageIds, queueMessage.MessageId)
	}

	requestContext.HandleHttpResponse(
		buildSuccessfulResponse(
			fmt.Sprintf(SuccessfullyBatchMessage, batchSize),
			map[string]any{"MessageIds": processedMessageIds}))
}
