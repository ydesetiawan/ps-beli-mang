package dto

type OrderEstimateRequest struct {
	UserLocation UserLocation `json:"userLocation"`
	Orders       []Order      `json:"orders"`
}

type UserLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
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

type OrderEstimateResponse struct {
	TotalPrice                     float64 `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int     `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           string  `json:"calculatedEstimateId"`
}
