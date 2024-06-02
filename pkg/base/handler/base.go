package handler

import (
	"fmt"
	"go/types"
	"net/http"
	"ps-beli-mang/pkg/errs"
	"ps-beli-mang/pkg/httphelper"
	"ps-beli-mang/pkg/httphelper/response"
	"ps-beli-mang/pkg/middleware"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

type HandlerFn func(c echo.Context) *response.WebResponse

type BaseHTTPHandler struct {
	Handlers interface{}
	DB       types.Nil
	Params   map[string]string
}

func (h *BaseHTTPHandler) RunAction(fn HandlerFn) echo.HandlerFunc {
	return h.HandlePanic(h.Execute(fn))
}

func (h *BaseHTTPHandler) Execute(fn HandlerFn) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			// return error if err
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					httphelper.WriteJSON(c.Response(), http.StatusInternalServerError,
						response.WebResponse{
							Status:  http.StatusInternalServerError,
							Message: http.StatusText(http.StatusInternalServerError),
						},
					)
					return
				}

				resultError := errs.ErrorAdvisor(err)
				httphelper.WriteJSON(c.Response(), resultError.Status,
					response.WebResponse{
						Status:  resultError.Status,
						Message: resultError.Message,
						Error:   resultError.Error,
						Data:    types.Interface{},
					},
				)
			}
		}()

		resp := fn(c)
		httpStatus := resp.Status
		if resp.RawData == nil {
			httphelper.WriteJSON(c.Response(), httpStatus,
				response.WebResponse{
					Token:      resp.Token,
					Status:     httpStatus,
					Message:    resp.Message,
					Data:       resp.Data,
					Pagination: resp.Pagination})
		} else {
			httphelper.WriteJSON(c.Response(), httpStatus, resp.RawData)
		}

		return nil
	}
}

func (h *BaseHTTPHandler) HandlePanic(fn echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				h.logPanicMessage(c.Request(), "CaptureLastPanic NEED TO FIX NOW", err)
				httphelper.WriteJSON(c.Response(), http.StatusInternalServerError, "Request unhalted unexpectedly, please contact administrator")
			}
		}()
		// Call the handler function and return its error if any.
		return fn(c)
	}
}

func (h *BaseHTTPHandler) logPanicMessage(r *http.Request, message string, err interface{}) {
	errStack := errs.StackAndFile(3)
	errInfo := fmt.Sprintf("\n SCM-production service \n* MUST FIX \U0001f4a3 \U0001f4a3 \U0001f4a3 "+
		"Panic Error: %v*", err)
	msg := fmt.Sprintf("%s\n\nStack trace: \n%s...", errInfo, errStack)

	fmt.Println("\nPANIC:", msg)
	src := "\n--- (Staging " + r.Host + ") ---\n"

	slog.ErrorCtx(r.Context(), message+src+msg, "attrs", errs.GetDefaultRequestFields(r))
}

func (h *BaseHTTPHandler) RunActionAuth(fn HandlerFn) echo.HandlerFunc {
	return h.HandlePanic(h.ExecuteAuth(h.Execute(fn)))
}

func (h *BaseHTTPHandler) ExecuteAuth(fn echo.HandlerFunc) echo.HandlerFunc {
	return middleware.JWTAuthMiddleware(fn)
}
