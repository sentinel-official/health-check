package types

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

type Response struct {
	Success bool        `json:"success"`
	Error   *Error      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func NewResponseError(code int, message string) *Response {
	return &Response{
		Success: false,
		Error:   NewError(code, message),
	}
}

func NewResponseResult(v interface{}) *Response {
	return &Response{
		Success: true,
		Result:  v,
	}
}
