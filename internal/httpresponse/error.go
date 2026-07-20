package httpresponse

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, message string, details string) Error {
	return Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}
