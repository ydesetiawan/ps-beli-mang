package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
)

type OrderService interface {
	OrderEstimate(ctx context.Context, request dto.OrderEstimateRequest) (dto.OrderEstimateResponse, error)
}
