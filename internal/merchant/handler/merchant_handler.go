package handler

import (
	"github.com/labstack/echo/v4"
	"ps-beli-mang/internal/merchant/service"
	"ps-beli-mang/pkg/httphelper/response"
)

type MerchantHandler struct {
	merchantService service.MerchantService
}

func NewMerchantHandler(merchantService service.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		merchantService: merchantService,
	}
}

func (h *MerchantHandler) GetMerchant(ctx echo.Context) *response.WebResponse {
	return nil
}

func (h *MerchantHandler) CreateMerchant(ctx echo.Context) *response.WebResponse {
	return nil
}

func (h *MerchantHandler) CreateMerchantItem(ctx echo.Context) *response.WebResponse {
	return nil
}

func (h *MerchantHandler) GetMerchantItem(ctx echo.Context) *response.WebResponse {
	return nil
}
