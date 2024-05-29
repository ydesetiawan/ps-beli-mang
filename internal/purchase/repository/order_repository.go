package repository

import (
	"golang.org/x/net/context"
	merchantModel "ps-beli-mang/internal/merchant/model"
)

type OrderRepository interface {
	GetMerchantItem(ctx context.Context, args []interface{}) ([]merchantModel.MerchantItem, error)
}
