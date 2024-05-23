package server

import (
	"fmt"
	"log/slog"
	imagehandler "ps-beli-mang/internal/image/handler"
	userhandler "ps-beli-mang/internal/user/handler"
	bhandler "ps-beli-mang/pkg/base/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	baseHandler  *bhandler.BaseHTTPHandler
	userHandler  *userhandler.UserHandler
	imageHandler *imagehandler.ImageHandler
	echo         *echo.Echo
	port         int
}

func NewServer(
	bHandler *bhandler.BaseHTTPHandler,
	userHandler *userhandler.UserHandler,
	imageHandler *imagehandler.ImageHandler,
	port int,
) Server {
	return Server{
		baseHandler:  bHandler,
		userHandler:  userHandler,
		imageHandler: imageHandler,
		echo:         echo.New(),
		port:         port,
	}
}

func (s *Server) Run() error {
	slog.Info(fmt.Sprintf("Starting HTTP server at :%d ...", s.port))
	e := echo.New()

	//e.Validator = &helpers.CustomValidator{Validator: validator.New()}
	//e.HTTPErrorHandler = helpers.ErrorHandler
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true}))

	s.setupRouter(e)

	return e.Start(fmt.Sprintf(":%d", s.port))
}
