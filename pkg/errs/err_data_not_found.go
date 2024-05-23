package errs

import (
	"fmt"
	"net/http"
)

type ErrDataNotFound struct {
	ErrBase
	DataId any
}

func (e ErrDataNotFound) Error() string {
	return fmt.Sprintf("%s. data_id: %v", e.msg, e.DataId)
}

func NewErrDataNotFound1(msg string) ErrDataNotFound {
	return ErrDataNotFound{
		ErrBase: ErrBase{
			msg:           msg,
			ErrorData:     ErrorData{},
			StatusCode:    http.StatusNotFound,
			StatusMessage: http.StatusText(http.StatusNotFound),
		},
		DataId: nil,
	}
}

func NewErrDataNotFound(msg string, dataId any, errData ErrorData) ErrDataNotFound {
	return ErrDataNotFound{
		ErrBase: ErrBase{
			msg:           msg,
			ErrorData:     errData,
			StatusCode:    http.StatusNotFound,
			StatusMessage: http.StatusText(http.StatusNotFound),
		},
		DataId: dataId,
	}
}

func NewErrDataNotFoundWithStackTrace(msg string, dataId any, errData ErrorData) error {
	return WrapErrorWithStackTrace(
		ErrDataNotFound{
			ErrBase: ErrBase{
				msg:           msg,
				ErrorData:     errData,
				StatusCode:    http.StatusNotFound,
				StatusMessage: http.StatusText(http.StatusNotFound),
			},
			DataId: dataId,
		},
	)
}
