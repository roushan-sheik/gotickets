package httpresponse

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, message string, detail string) Error {
	return Error{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}
