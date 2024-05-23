package errs

import "net/http"

type ErrBadRequest struct {
	ErrBase
}

func (e ErrBadRequest) Error() string {
	return e.msg
}

func NewErrBadRequest(msg string) ErrBadRequest {
	return ErrBadRequest{
		ErrBase: ErrBase{
			msg:           msg,
			StatusCode:    http.StatusBadRequest,
			StatusMessage: http.StatusText(http.StatusBadRequest),
		},
	}
}

func NewErrBadRequestWithStackTrace(msg string) error {
	return WrapErrorWithStackTrace(
		ErrBadRequest{
			ErrBase: ErrBase{
				msg:           msg,
				StatusCode:    http.StatusBadRequest,
				StatusMessage: http.StatusText(http.StatusBadRequest),
			},
		},
	)
}
