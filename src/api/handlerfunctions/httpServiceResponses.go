package handlerfunctions

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

func buildSuccessfulResponse(time int64, responseMessage string, payload map[string]any) HttpServiceResponse {
	return HttpServiceResponse{
		Response:           SuccessfulResponse,
		RequestStartedTime: time,
		Message:            responseMessage,
		Payload:            payload,
	}
}

func buildErrorResponse(time int64, responseMessage string, payload map[string]any) HttpServiceResponse {
	return HttpServiceResponse{
		Response:           ErrorResponse,
		RequestStartedTime: time,
		Message:            responseMessage,
		Payload:            payload,
	}
}
