package service

import (
	"context"
	"ps-beli-mang/internal/merchant/dto"
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/merchant/repository"
	orderRepository "ps-beli-mang/internal/purchase/repository"
)

type merchantService struct {
	merchantRepository repository.MerchantRepository
	orderRepository    orderRepository.OrderRepository
}

func NewMerchantServiceImpl(merchantRepository repository.MerchantRepository, orderRepository orderRepository.OrderRepository) MerchantService {
	return &merchantService{
		merchantRepository: merchantRepository,
		orderRepository:    orderRepository,
	}
}

func (s *merchantService) CreateMerchant(ctx context.Context, req *dto.MerchantDto) (id string, err error) {
	id, err = s.merchantRepository.CreateMerchant(ctx, req)
	if err == nil {
		s.orderRepository.ClearCache()
	}
	return id, err
}

func (s *merchantService) GetMerchants(ctx context.Context, req *dto.MerchantQuery) (merchants []model.Merchant, total int, err error) {
	return s.merchantRepository.GetMerchants(ctx, req)
}

func (s *merchantService) CreateMerchantItem(ctx context.Context, merchantId string, req *dto.MerchantItemDto) (id string, err error) {
	id, err = s.merchantRepository.CreateMerchantItem(ctx, merchantId, req)
	if err == nil {
		s.orderRepository.ClearCache()
	}
	return id, err
}

func (s *merchantService) GetMerchantItems(ctx context.Context, merchantId string, req *dto.MerchantItemQuery) (merchantItems []model.MerchantItem, total int, err error) {
	return s.merchantRepository.GetMerchantItems(ctx, merchantId, req)
}
