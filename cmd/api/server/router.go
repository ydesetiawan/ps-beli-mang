package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) setupRouter(e *echo.Echo) {
	v1 := e.Group("")
	v1.GET("/health", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Health Check OK")
	})
	v1.POST("/admin/register", s.baseHandler.RunAction(s.userHandler.RegisterAdmin))
	v1.POST("/admin/login", s.baseHandler.RunAction(s.userHandler.LoginAdmin))
	v1.POST("/user/register", s.baseHandler.RunAction(s.userHandler.RegisterUser))
	v1.POST("/user/login", s.baseHandler.RunAction(s.userHandler.LoginUser))

	v1.POST("/image", s.baseHandler.RunActionAuth(s.imageHandler.UploadImage))
}
