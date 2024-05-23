package errs

import "net/http"

type ErrForbidden struct {
	ErrBase
}

func (e ErrForbidden) Error() string {
	return e.msg
}

func NewErrForbidden(msg string) ErrForbidden {
	return ErrForbidden{
		ErrBase: ErrBase{
			msg:           msg,
			StatusCode:    http.StatusForbidden,
			StatusMessage: http.StatusText(http.StatusForbidden),
		},
	}
}

func NewErrForbiddenWithStackTrace(msg string) error {
	return WrapErrorWithStackTrace(
		ErrForbidden{
			ErrBase: ErrBase{
				msg:           msg,
				StatusCode:    http.StatusForbidden,
				StatusMessage: http.StatusText(http.StatusForbidden),
			},
		},
	)
}
