package dto

import (
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/pkg/helper"

	"github.com/go-playground/validator/v10"
)

type Location struct {
	Lat  float64 `json:"lat" validate:"required"`
	Long float64 `json:"long" validate:"required"`
}

type MerchantDto struct {
	Name             string   `json:"name" validate:"required,min=2,max=30"`
	MerchantCategory string   `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
	ImageUrl         string   `json:"imageUrl" validate:"required,validUrl"`
	Location         Location `json:"location" validate:"required"`
}

type MerchantItemDto struct {
	MerchantId      string  `param:"merchantId"`
	Name            string  `json:"name" validate:"required,min=2,max=30"`
	ProductCategory string  `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	Price           float64 `json:"price" validate:"required,min=1"`
	ImageUrl        string  `json:"imageUrl" validate:"required,validUrl"`
}

type MerchantQuery struct {
	ID               string                 `query:"merchantId"`
	Name             string                 `query:"name"`
	MerchantCategory model.MerchantCategory `query:"merchantCategory"`
	SortCreatedAt    model.SortType         `query:"createdAt"`
	Limit            int                    `query:"limit"`
	Offset           int                    `query:"offset"`
}

func (f *MerchantQuery) Validate() {
	if f.Limit == 0 {
		f.Limit = 5
	}
	if f.Offset == 0 {
		f.Offset = 0
	}
}

type MerchantItemQuery struct {
	MerchantID      string             `param:"merchantId"`
	ID              string             `query:"itemId"`
	Name            string             `query:"name"`
	ProductCategory model.ItemCategory `query:"productCategory"`
	SortCreatedAt   model.SortType     `query:"createdAt"`
	Limit           int                `query:"limit"`
	Offset          int                `query:"offset"`
}

func (f *MerchantItemQuery) Validate() {
	if f.Limit == 0 {
		f.Limit = 5
	}
	if f.Offset == 0 {
		f.Offset = 0
	}
}

func ValidateMerchantReq(req *MerchantDto) error {
	validate := validator.New()

	validate.RegisterValidation("validUrl", helper.ValidateURL)

	return validate.Struct(req)
}

func ValidateMerchantItemReq(req *MerchantItemDto) error {
	validate := validator.New()

	validate.RegisterValidation("validUrl", helper.ValidateURL)

	return validate.Struct(req)
}

type MerchantDtoResponse struct {
	MerchantId       string   `json:"merchantId"`
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageUrl         string   `json:"imageUrl"`
	Location         Location `json:"location"`
	CreatedAt        string   `json:"createdAt"`
}

type MerchantItemDtoResponse struct {
	ItemId          string  `json:"itemId"`
	Name            string  `json:"name"`
	ProductCategory string  `json:"productCategory"`
	ImageUrl        string  `json:"imageUrl"`
	Price           float64 `json:"price"`
	CreatedAt       string  `json:"createdAt"`
}
