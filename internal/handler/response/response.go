package response

type ApiResponse struct {
	Code    int         `json:"statusCode"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}
