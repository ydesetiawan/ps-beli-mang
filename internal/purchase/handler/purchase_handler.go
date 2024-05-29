package handler

import (
	"github.com/labstack/echo/v4"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/service"
	"ps-beli-mang/pkg/helper"
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
	var request = new(dto.OrderEstimateRequest)
	err := ctx.Bind(&request)
	helper.Panic400IfError(err)

	result, err := h.orderService.OrderEstimate(ctx.Request().Context(), request)
	helper.PanicIfError(err, "CreateMedicalPatient failed")

	return &response.WebResponse{
		Status:  201,
		Message: "OrderEstimate Successfully",
		Data:    result,
	}
}

func (h *PurchaseHandler) Order(ctx echo.Context) *response.WebResponse {
	return nil
}

func (h *PurchaseHandler) GetOrder(ctx echo.Context) *response.WebResponse {
	return nil
}
