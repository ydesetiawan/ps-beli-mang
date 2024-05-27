package handler

import (
	"github.com/labstack/echo/v4"
	"ps-beli-mang/internal/purchase/service"
	"ps-beli-mang/pkg/httphelper/response"
)

type PurchaseHandler struct {
	orderService service.OrderService
}

func NewPurchaseHandler(orderService service.OrderService) *PurchaseHandler {
	return &PurchaseHandler{
		orderService: orderService,
	}
}

func (h *PurchaseHandler) GetNearMerchant(ctx echo.Context) *response.WebResponse {
	return nil
}

func (h *PurchaseHandler) OrderEstimate(ctx echo.Context) *response.WebResponse {
	return nil
}

func (h *PurchaseHandler) Order(ctx echo.Context) *response.WebResponse {
	return nil
}

func (h *PurchaseHandler) GetOrder(ctx echo.Context) *response.WebResponse {
	return nil
}
