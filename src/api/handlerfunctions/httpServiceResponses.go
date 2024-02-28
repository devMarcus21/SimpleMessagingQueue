package handlerfunctions

import (
	"github.com/google/uuid"
)

const (
	SuccessfulResponse        string = "Success"
	SuccessfulResponseMessage string = "Success"
	ErrorResponse             string = "Error"

	SuccessfullyBatchMessage string = "Successfully added batch of size: %d"
	QueueIsEmptyMessage      string = "Queue is empty"
)

type HttpServiceResponse struct {
	Response           string
	RequestId          uuid.UUID
	RequestStartedTime int64
	Message            string
	Payload            map[string]any
}

func buildSuccessfulResponse(responseMessage string, payload map[string]any) HttpServiceResponse {
	return HttpServiceResponse{
		Response: SuccessfulResponse,
		Message:  responseMessage,
		Payload:  payload,
	}
}

func buildErrorResponse(responseMessage string, payload map[string]any) HttpServiceResponse {
	return HttpServiceResponse{
		Response: ErrorResponse,
		Message:  responseMessage,
		Payload:  payload,
	}
}
