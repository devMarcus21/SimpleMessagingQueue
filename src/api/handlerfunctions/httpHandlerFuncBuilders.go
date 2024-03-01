package handlerfunctions

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	asyncQueueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/asyncqueue"
	configuration "github.com/devMarcus21/SimpleMessagingQueue/src/configuration"
	logging "github.com/devMarcus21/SimpleMessagingQueue/src/logging"

	"github.com/google/uuid"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type HandlerRequestContext struct {
	httpWriter               http.ResponseWriter
	httpReader               *http.Request
	logger                   *slog.Logger
	requestId                uuid.UUID
	environmentConfiguration configuration.Configuration
	epochRequestStartTime    int64
}

func (context *HandlerRequestContext) Logger() *slog.Logger {
	return context.logger
}

func (context *HandlerRequestContext) GetHttpBody() io.ReadCloser {
	return context.httpReader.Body
}

func (context *HandlerRequestContext) RequestId() uuid.UUID {
	return context.requestId
}

func (context *HandlerRequestContext) RequestStartTime() int64 {
	return context.epochRequestStartTime
}

func (context *HandlerRequestContext) Configuration() configuration.Configuration {
	return context.environmentConfiguration
}

func (context *HandlerRequestContext) AddHttpStatusCode(httpStatusCode int) {
	context.httpWriter.WriteHeader(httpStatusCode)
}

func (context *HandlerRequestContext) HandleHttpResponse(response HttpServiceResponse) {
	response.RequestId = context.requestId
	response.RequestStartedTime = context.epochRequestStartTime

	json.NewEncoder(context.httpWriter).Encode(response)
}

func BuildHttpHandlerFunc(requestHandler func(HandlerRequestContext, asyncQueueUtils.AsyncQueueWrapper), loggerBuilder logging.LoggerBuilder, asyncQueue asyncQueueUtils.AsyncQueueWrapper, config configuration.Configuration, handlerActionName logging.LogName) HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		requestId := uuid.New()
		logger := loggerBuilder(logging.Request).With( // Adds the given logging attributes to every single call to the logger
			// Key Value
			"RequestId", requestId,
			logging.HandlerActionName.String(), handlerActionName.String())

		epochTimeNow := time.Now().Unix()

		requestContext := HandlerRequestContext{
			httpWriter:               writer,
			httpReader:               reader,
			logger:                   logger,
			requestId:                requestId,
			environmentConfiguration: config,
			epochRequestStartTime:    epochTimeNow,
		}

		requestHandler(requestContext, asyncQueue)
	}
}
