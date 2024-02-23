package errorresponses

type ApiErrorResponses int64

const (
	ApiErrorIota string = "ApiErrorIota"

	GivenEmptyBatchError ApiErrorResponses = iota
	BatchSizeBiggerThanMaxBatchSizeError
	JsonUnmarshalError
)

var apiErrorResponseNames = map[ApiErrorResponses]string{
	GivenEmptyBatchError:                 "GivenEmptyBatchError",
	BatchSizeBiggerThanMaxBatchSizeError: "BatchSizeBiggerThanMaxBatchSizeError",
	JsonUnmarshalError:                   "JsonUnmarshalError",
}

var apiErrorResponseMessages = map[ApiErrorResponses]string{
	GivenEmptyBatchError:                 "Error: received empty batch",
	BatchSizeBiggerThanMaxBatchSizeError: "Error: batch size (%d) is greater than max size allowed (%d)",
	JsonUnmarshalError:                   "Error: failed to decode request body",
}

func (err ApiErrorResponses) String() string {
	return apiErrorResponseNames[err]
}

func (err ApiErrorResponses) Message() string {
	return apiErrorResponseMessages[err]
}
