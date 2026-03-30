package dto

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message"`
}

func SuccessResp(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Message: "success",
	}
}

func ErrorResp(err string, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
		Message: message,
	}
}
