package model

type ApiResponse[T any] struct {
	Data          T        `json:"data,omitempty"`
	Message       string   `json:"message,omitempty"`
	Success       bool     `json:"success"`
	Errors        []string `json:"errors,omitempty"`
	SecondaryData any      `json:"secondary_data,omitempty"`
}

func NewApiResponse[T any](data T, message string) ApiResponse[T] {
	apiResponse := ApiResponse[T]{
		Data:    data,
		Success: true,
		Message: message,
	}
	return apiResponse
}

func NewApiResponseMessage[T any](message string, success bool) ApiResponse[T] {
	apiResponse := ApiResponse[T]{
		Success: success,
		Message: message,
	}
	return apiResponse
}

func NewApiResponseWithSecondary[T any](data T, message string, secondary any) ApiResponse[T] {
	apiResponse := ApiResponse[T]{
		Data:          data,
		Success:       true,
		Message:       message,
		SecondaryData: secondary,
	}
	return apiResponse
}

func NewApiResponseError(errors []string, message string) ApiResponse[any] {
	if errors == nil && message != "" {
		errors = []string{message}
	}
	if errors != nil && message == "" {
		message = errors[0]
	}
	apiResponse := ApiResponse[any]{
		Success: false,
		Errors:  errors,
		Message: message,
	}
	return apiResponse
}

func (a *ApiResponse[T]) SetSecondaryData(secondary any) *ApiResponse[T] {
	a.SecondaryData = secondary
	return a
}

func BuildSecondary[T any](data T, message string) ApiResponse[any] {
	apiResponse := ApiResponse[any]{
		Data:    data,
		Success: true,
		Message: message,
	}
	return apiResponse
}
