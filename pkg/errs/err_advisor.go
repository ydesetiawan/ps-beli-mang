package errs

import (
	"errors"
	"net/http"
	"ps-beli-mang/pkg/httphelper/response"
)

func ErrorAdvisor(err error) response.WebResponse {
	errBase := NewErrBase(http.StatusInternalServerError, "please try again", ErrorData{Process: InternalServer, DataId: ""})
	resp := response.WebResponse{}
	resp.Status = errBase.StatusCode
	resp.Message = errBase.StatusMessage
	resp.Error = errBase.ErrorData
	resp.Data = err.Error()

	var errBadRequest ErrBadRequest
	var errDataConflict ErrDataConflict
	var errDataNotFound ErrDataNotFound
	var errForbidden ErrForbidden
	var errUnauthorized ErrUnauthorized
	var errUnprocessableEntity ErrUnprocessableEntity
	switch {
	case errors.As(err, &errUnprocessableEntity):
		resp.Status = errUnprocessableEntity.StatusCode
		resp.Message = errUnprocessableEntity.Error()
		resp.Error = errUnprocessableEntity.ErrorData
	case errors.As(err, &errBadRequest):
		resp.Status = errBadRequest.StatusCode
		resp.Message = errBadRequest.Error()
		resp.Error = errBadRequest.ErrorData
	case errors.As(err, &errDataConflict):
		resp.Status = errDataConflict.StatusCode
		resp.Message = errDataConflict.Error()
		resp.Error = errDataConflict.ErrorData
	case errors.As(err, &errDataNotFound):
		resp.Status = errDataNotFound.StatusCode
		resp.Message = errDataNotFound.Error()
		resp.Error = errDataNotFound.ErrorData
	case errors.As(err, &errForbidden):
		resp.Status = errForbidden.StatusCode
		resp.Message = errForbidden.Error()
		resp.Error = errForbidden.ErrorData
	case errors.As(err, &errUnauthorized):
		resp.Status = errUnauthorized.StatusCode
		resp.Message = errUnauthorized.Error()
		resp.Error = errUnauthorized.ErrorData
	}

	return resp
}
