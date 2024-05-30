package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
)

func (o orderService) CreateOrder(ctx context.Context, request dto.CreateOrderRequest) error {
	return o.orderRepository.UpdateOrderSetIsOrderTrue(ctx, request.CalculatedEstimateId)
}
