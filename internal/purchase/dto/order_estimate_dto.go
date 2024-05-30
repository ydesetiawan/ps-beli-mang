package dto

import (
	"ps-beli-mang/internal/merchant/model"
)

type OrderEstimateRequest struct {
	UserLocation model.Location `json:"userLocation"`
	Orders       []Order        `json:"orders"`
	UserID       string         `json:"-"`
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

type PreOrder struct {
	MerchantItems           []model.MerchantItem
	ItemQtyIds              map[string]int
	MerchantStartingPointId string
	UserLocation            model.Location
	UserID                  string
}

type OrderEstimateResponse struct {
	TotalPrice                     float64 `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int     `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           string  `json:"calculatedEstimateId"`
}

type CreateOrderRequest struct {
	CalculatedEstimateId string
}

type CreateOrderResponse struct {
	OrderId string
}
