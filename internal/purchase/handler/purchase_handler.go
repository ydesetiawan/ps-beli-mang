package handler

import (
	"github.com/labstack/echo/v4"
	merchantModel "ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/service"
	"ps-beli-mang/internal/user/model"
	userService "ps-beli-mang/internal/user/service"
	"ps-beli-mang/pkg/base/handler"
	"ps-beli-mang/pkg/errs"
	"ps-beli-mang/pkg/helper"
	"ps-beli-mang/pkg/httphelper/response"
	"strconv"
	"strings"
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

func getUserLocation(c echo.Context) (merchantModel.Location, error) {
	defaultResponse := merchantModel.Location{}
	locArr := strings.Split(c.ParamValues()[0], ",")
	if len(locArr) != 2 {
		return defaultResponse, errs.NewErrBadRequest("Invalid latitude, longitude")
	}
	latStr := locArr[0]
	longStr := locArr[1]

	// Convert latitude and longitude strings to float64
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return defaultResponse, errs.NewErrBadRequest("Invalid latitude")
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return defaultResponse, errs.NewErrBadRequest("Invalid longitude")
	}

	return merchantModel.Location{Lat: lat, Long: long}, nil
}

func (h *PurchaseHandler) GetNearbyMerchant(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleUser(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	var params = new(dto.MerchantRequestParams)
	err = ctx.Bind(params)
	helper.Panic400IfError(err)

	params.Validate()

	params.UserLocation, err = getUserLocation(ctx)
	helper.Panic400IfError(err)

	result, total, err := h.orderService.GetNearbyMerchants(ctx.Request().Context(), *params)
	helper.PanicIfError(err, "GetNearbyMerchant failed")

	type Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	return &response.WebResponse{
		Status:  200,
		Message: "GetNearbyMerchant Successfully",
		RawData: struct {
			Data []dto.GetNearbyMerchantResponse `json:"data"`
			Meta Meta                            `json:"meta"`
		}{
			Data: result,
			Meta: Meta{
				Limit:  params.Limit,
				Offset: params.Offset,
				Total:  total,
			},
		},
	}
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

	var params = new(dto.MerchantRequestParams)
	err = ctx.Bind(params)
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
