package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/repository"
)

type OrderService interface {
	OrderEstimate(ctx context.Context, request dto.OrderEstimateRequest) (dto.OrderEstimateResponse, error)
	CreateOrder(ctx context.Context, request dto.CreateOrderRequest) error
}

type orderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderServiceImpl(orderRepository repository.OrderRepository) OrderService {
	return &orderService{
		orderRepository: orderRepository,
	}
}
