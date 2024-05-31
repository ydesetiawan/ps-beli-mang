package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
)

func (o orderService) GetNearbyMerchants(ctx context.Context, params dto.MerchantRequestParams) ([]dto.GetNearbyMerchantResponse, error) {
	return o.orderRepository.GetAllMerchants(ctx)
}
