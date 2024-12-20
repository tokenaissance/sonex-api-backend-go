package response

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: "Success",
		Data:    data,
	}
}

func Error(code int, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
	}
}
