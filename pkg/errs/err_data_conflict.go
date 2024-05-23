package errs

import (
	"fmt"
	"net/http"
)

type ErrDataConflict struct {
	ErrBase
	DataId any
}

func (e ErrDataConflict) Error() string {
	return fmt.Sprintf("%s %v", e.msg, e.DataId)
}

func NewErrDataConflict(msg string, dataId any) ErrDataConflict {
	return ErrDataConflict{
		ErrBase: ErrBase{
			msg:           msg,
			StatusCode:    http.StatusConflict,
			StatusMessage: http.StatusText(http.StatusConflict),
		},
		DataId: dataId,
	}
}

func NewErrDataConflictWithStackTrace(msg string, dataId any) error {
	return WrapErrorWithStackTrace(
		ErrDataConflict{
			ErrBase: ErrBase{
				msg:           msg,
				StatusCode:    http.StatusConflict,
				StatusMessage: http.StatusText(http.StatusConflict),
			},
			DataId: dataId,
		},
	)
}
