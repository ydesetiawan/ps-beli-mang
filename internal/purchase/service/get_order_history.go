package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
)

func (o orderService) GetOrders(ctx context.Context, params dto.MerchantRequestParams) ([]dto.OrderDataResponse, error) {
	return o.orderRepository.GetOrdersByUser(ctx, params)
}
