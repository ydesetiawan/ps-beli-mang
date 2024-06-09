package dto

import (
	"github.com/go-playground/validator/v10"
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

func MergeOrders(request *OrderEstimateRequest) OrderEstimateRequest {
	mergedOrdersMap := make(map[string]Order)

	for _, order := range request.Orders {
		if existingOrder, found := mergedOrdersMap[order.MerchantID]; found {
			// If the order already exists, merge the items
			for _, newItem := range order.Items {
				found := false
				for i, existingItem := range existingOrder.Items {
					if existingItem.ItemID == newItem.ItemID {
						existingOrder.Items[i].Quantity += newItem.Quantity
						found = true
						break
					}
				}
				if !found {
					existingOrder.Items = append(existingOrder.Items, newItem)
				}
			}
			mergedOrdersMap[order.MerchantID] = existingOrder
		} else {
			// If the order does not exist, add it to the map
			mergedOrdersMap[order.MerchantID] = order
		}
	}

	// Convert the map back to a slice
	var mergedOrders []Order
	for _, order := range mergedOrdersMap {
		mergedOrders = append(mergedOrders, order)
	}

	return OrderEstimateRequest{
		UserLocation: request.UserLocation,
		Orders:       mergedOrders,
	}
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
	CalculatedEstimateId string `json:"calculatedEstimateId" validate:"required"`
}

func ValidateCreateOrderRequest(req *CreateOrderRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

type CreateOrderResponse struct {
	OrderId string `json:"orderId"`
}
