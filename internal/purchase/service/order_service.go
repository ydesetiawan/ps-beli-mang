package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/repository"
)

type OrderService interface {
	GetNearbyMerchants(ctx context.Context, params dto.MerchantRequestParams) (result []dto.GetNearbyMerchantResponse, total int, err error)
	OrderEstimate(ctx context.Context, request dto.OrderEstimateRequest) (dto.OrderEstimateResponse, error)
	CreateOrder(ctx context.Context, request dto.CreateOrderRequest) error
	GetOrders(ctx context.Context, params dto.MerchantRequestParams) ([]dto.OrderDataResponse, error)
}

type orderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderServiceImpl(orderRepository repository.OrderRepository) OrderService {
	return &orderService{
		orderRepository: orderRepository,
	}
}
