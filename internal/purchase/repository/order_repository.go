package repository

import (
	"golang.org/x/net/context"
	merchantModel "ps-beli-mang/internal/merchant/model"
)

type OrderRepository interface {
	GetMerchantItems(ctx context.Context, args []interface{}) ([]merchantModel.MerchantItem, error)
}
