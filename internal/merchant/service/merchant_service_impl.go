package service

import (
	"ps-beli-mang/internal/merchant/repository"
)

type merchantService struct {
	merchantRepository repository.MerchantRepository
}

func NewMerchantServiceImpl(merchantRepository repository.MerchantRepository) MerchantService {
	return &merchantService{
		merchantRepository: merchantRepository,
	}
}
