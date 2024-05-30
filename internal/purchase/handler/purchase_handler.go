package handler

import (
	"github.com/labstack/echo/v4"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/service"
	"ps-beli-mang/internal/user/model"
	userService "ps-beli-mang/internal/user/service"
	"ps-beli-mang/pkg/base/handler"
	"ps-beli-mang/pkg/helper"
	"ps-beli-mang/pkg/httphelper/response"
)

type PurchaseHandler struct {
	orderService service.OrderService
	userService  userService.UserService
}

func NewPurchaseHandler(orderService service.OrderService, userService userService.UserService) *PurchaseHandler {
	return &PurchaseHandler{
		orderService: orderService,
		userService:  userService,
	}
}

func hasAuthorizeRoleUser(ctx echo.Context, h *PurchaseHandler) (string, error) {
	userId, err := handler.GetUserId(ctx)
	if err != nil {
		return "", err
	}

	_, err = h.userService.GetUserByIdAndRole(ctx.Request().Context(), userId, string(model.USER))
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (h *PurchaseHandler) GetNearMerchant(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleUser(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	return nil
}

func (h *PurchaseHandler) OrderEstimate(ctx echo.Context) *response.WebResponse {
	userID, err := hasAuthorizeRoleUser(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	var request = new(dto.OrderEstimateRequest)
	err = ctx.Bind(&request)
	helper.Panic400IfError(err)

	request.UserID = userID
	result, err := h.orderService.OrderEstimate(ctx.Request().Context(), *request)
	helper.PanicIfError(err, "OrderEstimate failed")

	return &response.WebResponse{
		Status:  201,
		Message: "OrderEstimate Successfully",
		Data:    result,
	}
}

func (h *PurchaseHandler) Order(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleUser(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	var request = new(dto.CreateOrderRequest)
	err = ctx.Bind(&request)
	helper.Panic400IfError(err)

	err = h.orderService.CreateOrder(ctx.Request().Context(), *request)
	helper.PanicIfError(err, "CreateMedicalPatient failed")

	return &response.WebResponse{
		Status:  201,
		Message: "Order Successfully",
		Data: &dto.CreateOrderResponse{
			OrderId: request.CalculatedEstimateId,
		},
	}
}

func (h *PurchaseHandler) GetOrders(ctx echo.Context) *response.WebResponse {
	userID, err := hasAuthorizeRoleUser(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	var params = new(dto.OrderDataRequestParams)
	err = ctx.Bind(&params)
	helper.Panic400IfError(err)

	params.UserID = userID
	result, err := h.orderService.GetOrders(ctx.Request().Context(), *params)
	helper.PanicIfError(err, "GetOrders failed")

	return &response.WebResponse{
		Status:  200,
		Message: "GetOrders Successfully",
		Data:    result,
	}
}
