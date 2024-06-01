package repository

import (
	merchantModel "ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/model"

	"golang.org/x/net/context"
)

type OrderRepository interface {
	GetMerchantItems(ctx context.Context, args []interface{}) ([]merchantModel.MerchantItem, error)
	SaveOrder(ctx context.Context, order model.Order) error
	UpdateOrderSetIsOrderTrue(ctx context.Context, orderID string) error
	GetNearbyMerchantByUser(ctx context.Context, params dto.MerchantRequestParams) ([]dto.GetNearbyMerchantResponse, error)
	GetOrdersByUser(ctx context.Context, params dto.MerchantRequestParams) ([]dto.OrderDataResponse, error)
	GetAllMerchants(ctx context.Context) ([]dto.GetNearbyMerchantResponse, error)
	ClearCache()
}
