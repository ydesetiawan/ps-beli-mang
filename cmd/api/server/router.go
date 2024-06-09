package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) setupRouter(e *echo.Echo) {
	adminV1 := e.Group("/admin")
	adminV1.POST("/register", s.baseHandler.RunAction(s.userHandler.RegisterAdmin))
	adminV1.POST("/login", s.baseHandler.RunAction(s.userHandler.LoginAdmin))
	adminV1.POST("/merchants", s.baseHandler.RunActionAuth(s.merchantHandler.CreateMerchant))
	adminV1.GET("/merchants", s.baseHandler.RunActionAuth(s.merchantHandler.GetMerchant))
	adminV1.POST("/merchants/:merchantId/items", s.baseHandler.RunActionAuth(s.merchantHandler.CreateMerchantItem))
	adminV1.GET("/merchants/:merchantId/items", s.baseHandler.RunActionAuth(s.merchantHandler.GetMerchantItem))

	usersV1 := e.Group("/users")
	usersV1.POST("/login", s.baseHandler.RunAction(s.userHandler.LoginUser))
	usersV1.POST("/register", s.baseHandler.RunAction(s.userHandler.RegisterUser))
	usersV1.POST("/estimate", s.baseHandler.RunActionAuth(s.purchaseHandler.OrderEstimate))
	usersV1.POST("/orders", s.baseHandler.RunActionAuth(s.purchaseHandler.Order))
	usersV1.GET("/orders", s.baseHandler.RunActionAuth(s.purchaseHandler.GetOrders))

	e.POST("/image", s.baseHandler.RunActionAuth(s.imageHandler.UploadImage))
	e.GET("/merchants/nearby/:lat,:long", s.baseHandler.RunActionAuth(s.purchaseHandler.GetNearbyMerchant))
	e.GET("/health", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Health Check OK New 9 jun 2024 16:19")
	})
}
