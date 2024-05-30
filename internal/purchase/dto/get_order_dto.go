package dto

import (
	"ps-beli-mang/internal/merchant/model"
	"time"
)

type Merchant struct {
	MerchantID       string         `json:"merchantId"`
	Name             string         `json:"name"`
	MerchantCategory string         `json:"merchantCategory"`
	ImageURL         string         `json:"imageUrl"`
	Location         model.Location `json:"location"`
	CreatedAt        time.Time      `json:"createdAt"` // use time.Time to parse ISO 8601 format
}

type PurchaseItem struct {
	ItemID          string    `json:"itemId"`
	Name            string    `json:"name"`
	ProductCategory string    `json:"productCategory"`
	Price           float64   `json:"price"`
	Quantity        int       `json:"quantity"`
	ImageURL        string    `json:"imageUrl"`
	CreatedAt       time.Time `json:"createdAt"` // use time.Time to parse ISO 8601 format
}

type PurchaseOrder struct {
	Merchant Merchant       `json:"merchant"`
	Items    []PurchaseItem `json:"items"`
}

type OrderDataResponse struct {
	OrderID string          `json:"orderId"`
	Orders  []PurchaseOrder `json:"orders"`
}

type OrderDataRequestParams struct {
	MerchantID       string `query:"merchantId"`
	Limit            int    `query:"limit"`
	Offset           int    `query:"offset"`
	Name             string `query:"name"`
	MerchantCategory string `query:"merchantCategory"`
	UserID           string `query:"-"`
}
