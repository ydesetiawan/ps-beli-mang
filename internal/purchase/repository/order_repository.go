package repository

import (
	"golang.org/x/net/context"
	merchantModel "ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/model"
)

type OrderRepository interface {
	GetMerchantItems(ctx context.Context, args []interface{}) ([]merchantModel.MerchantItem, error)
	SaveOrder(ctx context.Context, order model.Order) error
	UpdateIsOrderTrue(ctx context.Context, orderID string) error
}
