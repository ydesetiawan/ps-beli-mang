package handler

import (
	"ps-beli-mang/internal/merchant/dto"
	"ps-beli-mang/internal/merchant/service"
	"ps-beli-mang/internal/user/model"
	userService "ps-beli-mang/internal/user/service"
	"ps-beli-mang/pkg/base/handler"
	"ps-beli-mang/pkg/helper"
	"ps-beli-mang/pkg/httphelper/response"
	"time"

	"github.com/labstack/echo/v4"
)

type MerchantHandler struct {
	merchantService service.MerchantService
	userService     userService.UserService
}

func NewMerchantHandler(merchantService service.MerchantService, userService userService.UserService) *MerchantHandler {
	return &MerchantHandler{
		merchantService: merchantService,
		userService:     userService,
	}
}

func hasAuthorizeRoleAdmin(ctx echo.Context, h *MerchantHandler) (string, error) {
	userId, err := handler.GetUserId(ctx)
	if err != nil {
		return "", err
	}

	_, err = h.userService.GetUserByIdAndRole(ctx.Request().Context(), userId, string(model.ADMIN))
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (h *MerchantHandler) GetMerchant(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleAdmin(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	filter := dto.MerchantQuery{}
	err = ctx.Bind(&filter)
	helper.Panic400IfError(err)

	filter.Validate()

	merchants, total, err := h.merchantService.GetMerchants(ctx.Request().Context(), &filter)
	helper.PanicIfError(err, "GetMerchant failed")

	responseData := []dto.MerchantDtoResponse{}
	for _, merchant := range merchants {
		responseData = append(responseData, dto.MerchantDtoResponse{
			MerchantId:       merchant.ID,
			Name:             merchant.Name,
			MerchantCategory: merchant.Category.String(),
			ImageUrl:         merchant.ImageUrl,
			Location:         dto.Location{Lat: merchant.LocLat, Long: merchant.LocLong},
			CreatedAt:        merchant.CreatedAt.Format(time.RFC3339Nano),
		})
	}

	type Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	return &response.WebResponse{
		Status: 200,
		RawData: struct {
			Data []dto.MerchantDtoResponse `json:"data"`
			Meta Meta                      `json:"meta"`
		}{
			Data: responseData,
			Meta: Meta{
				Limit:  filter.Limit,
				Offset: filter.Offset,
				Total:  total,
			},
		},
	}
}

func (h *MerchantHandler) CreateMerchant(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleAdmin(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	request := dto.MerchantDto{}
	err = ctx.Bind(&request)
	helper.Panic400IfError(err)

	err = dto.ValidateMerchantReq(&request)
	helper.Panic400IfError(err)

	id, err := h.merchantService.CreateMerchant(ctx.Request().Context(), &request)
	helper.PanicIfError(err, "CreateMerchant failed")

	return &response.WebResponse{
		Status: 201,
		RawData: struct {
			MerchantId string `json:"merchantId"`
		}{MerchantId: id},
	}
}

func (h *MerchantHandler) CreateMerchantItem(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleAdmin(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	request := dto.MerchantItemDto{}
	err = ctx.Bind(&request)
	helper.Panic400IfError(err)

	err = dto.ValidateMerchantItemReq(&request)
	helper.Panic400IfError(err)

	id, err := h.merchantService.CreateMerchantItem(ctx.Request().Context(), request.MerchantId, &request)
	helper.PanicIfError(err, "CreateMerchantItem failed")

	return &response.WebResponse{
		Status: 201,
		RawData: struct {
			ItemId string `json:"itemId"`
		}{ItemId: id},
	}
}

func (h *MerchantHandler) GetMerchantItem(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleAdmin(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	filter := dto.MerchantItemQuery{}
	err = ctx.Bind(&filter)
	helper.Panic400IfError(err)

	filter.Validate()

	merchantItems, total, err := h.merchantService.GetMerchantItems(ctx.Request().Context(), filter.MerchantID, &filter)
	helper.PanicIfError(err, "GetMerchantItem failed")

	responseData := []dto.MerchantItemDtoResponse{}
	for _, merchantItem := range merchantItems {
		responseData = append(responseData, dto.MerchantItemDtoResponse{
			ItemId:          merchantItem.ID,
			Name:            merchantItem.Name,
			ProductCategory: merchantItem.Category.String(),
			ImageUrl:        merchantItem.ImageUrl,
			Price:           merchantItem.Price,
			CreatedAt:       merchantItem.CreatedAt.Format(time.RFC3339Nano),
		})
	}

	type Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	return &response.WebResponse{
		Status: 200,
		RawData: struct {
			Data []dto.MerchantItemDtoResponse `json:"data"`
			Meta Meta                          `json:"meta"`
		}{
			Data: responseData,
			Meta: Meta{
				Limit:  filter.Limit,
				Offset: filter.Offset,
				Total:  total,
			},
		},
	}
}
