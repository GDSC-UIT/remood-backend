package models

type Response struct {
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Data    interface{} `json:"data,omitempty"`
}


func ErrorResponse(message string) Response {
	return Response{
		Message: message,
		Error: true,
	}
}
 
func SuccessResponse(message string, data any) Response {
	return Response{
		Message: message,
		Error: false,
		Data: data,
	}
}