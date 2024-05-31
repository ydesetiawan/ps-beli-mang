package repository

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
)

func (o orderRepositoryImpl) GetNearbyMerchantByUser(ctx context.Context, params dto.MerchantRequestParams) ([]dto.GetNearbyMerchantResponse, error) {
	return o.GetAllMerchants(ctx)
}
