package errs

import "net/http"

type ErrBase struct {
	StatusCode    int
	StatusMessage string
	ErrorData     ErrorData
	msg           string
}

type ErrorData struct {
	Process string
	DataId  any
}

func (e ErrBase) Error() string {
	return e.msg
}

func NewErrBase(statusCode int, msg string, errData ErrorData) ErrBase {
	return ErrBase{
		StatusCode:    statusCode,
		StatusMessage: http.StatusText(statusCode),
		ErrorData:     errData,
		msg:           msg,
	}
}

func NewErrBaseWithStackTrace(statusCode int, msg string, errData ErrorData) error {
	return WrapErrorWithStackTrace(
		ErrBase{
			StatusCode:    statusCode,
			StatusMessage: http.StatusText(statusCode),
			ErrorData:     errData,
			msg:           msg,
		},
	)
}
