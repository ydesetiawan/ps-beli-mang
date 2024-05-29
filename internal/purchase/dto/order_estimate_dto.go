package dto

import (
	"ps-beli-mang/internal/merchant/model"
)

type OrderEstimateRequest struct {
	UserLocation model.Location `json:"userLocation"`
	Orders       []Order        `json:"orders"`
}

type Order struct {
	MerchantID      string `json:"merchantId"`
	IsStartingPoint bool   `json:"isStartingPoint"`
	Items           []Item `json:"items"`
}

type Item struct {
	ItemID   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

type OrderEstimateProcess struct {
	MerchantItems           []model.MerchantItem
	ItemQtyIds              map[string]int
	MerchantStartingPointId string
}

type OrderEstimateResponse struct {
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           string `json:"calculatedEstimateId"`
}
