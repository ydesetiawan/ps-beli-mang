package errs

import (
	"fmt"
	"net/http"
)

type ErrUnprocessableEntity struct {
	ErrBase
	DataId any
}

func (e ErrUnprocessableEntity) Error() string {
	return fmt.Sprintf("%s %v", e.msg, e.DataId)
}

func NewErrUnprocessableEntity(msg string, dataId any, errData ErrorData) ErrUnprocessableEntity {
	return ErrUnprocessableEntity{
		ErrBase: ErrBase{
			msg:           msg,
			StatusCode:    http.StatusUnprocessableEntity,
			ErrorData:     errData,
			StatusMessage: http.StatusText(http.StatusUnprocessableEntity),
		},
		DataId: dataId,
	}

}

func NewErrUnprocessableEntityWithStackTrace(msg string, dataId any, errData ErrorData) error {
	return WrapErrorWithStackTrace(
		ErrUnprocessableEntity{
			ErrBase: ErrBase{
				msg:           msg,
				ErrorData:     errData,
				StatusCode:    http.StatusUnprocessableEntity,
				StatusMessage: http.StatusText(http.StatusUnprocessableEntity),
			},
			DataId: dataId,
		},
	)
}

func NewErrInternalServerError(msg string, dataId any, errData ErrorData) error {
	return WrapErrorWithStackTrace(
		ErrUnprocessableEntity{
			ErrBase: ErrBase{
				msg:           msg,
				ErrorData:     errData,
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: http.StatusText(http.StatusInternalServerError),
			},
			DataId: dataId,
		},
	)
}

func NewErrInternalServerErrors(msg string, dataId any) error {
	return WrapErrorWithStackTrace(
		ErrUnprocessableEntity{
			ErrBase: ErrBase{
				msg:           msg,
				ErrorData:     ErrorData{},
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: http.StatusText(http.StatusInternalServerError),
			},
			DataId: dataId,
		},
	)
}
