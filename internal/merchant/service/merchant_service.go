package service

import (
	"context"
	"ps-beli-mang/internal/merchant/dto"
	"ps-beli-mang/internal/merchant/model"
)

type MerchantService interface {
	CreateMerchant(ctx context.Context, req *dto.MerchantDto) (id string, err error)
	GetMerchants(ctx context.Context, req *dto.MerchantQuery) (merchants []model.Merchant, total int, err error)
	CreateMerchantItem(ctx context.Context, merchantId string, req *dto.MerchantItemDto) (id string, err error)
	GetMerchantItems(ctx context.Context, merchantId string, req *dto.MerchantItemQuery) (merchantItems []model.MerchantItem, total int, err error)
}
