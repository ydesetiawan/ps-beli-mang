package errs

import "net/http"

type ErrUnauthorized struct {
	ErrBase
}

func (e ErrUnauthorized) Error() string {
	return e.msg
}

func NewErrUnauthorized(msg string) ErrUnauthorized {
	return ErrUnauthorized{
		ErrBase: ErrBase{
			msg:           msg,
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: http.StatusText(http.StatusUnauthorized),
		},
	}
}

func NewErrUnauthorizedWithStackTrace(msg string) error {
	return WrapErrorWithStackTrace(
		ErrUnauthorized{
			ErrBase: ErrBase{
				msg:           msg,
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: http.StatusText(http.StatusUnauthorized),
			},
		},
	)
}
