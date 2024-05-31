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

type MerchantItem struct {
	ItemID          string    `json:"itemId"`
	Name            string    `json:"name"`
	ProductCategory string    `json:"productCategory"`
	Price           float64   `json:"price"`
	Quantity        int       `json:"quantity"`
	ImageURL        string    `json:"imageUrl"`
	CreatedAt       time.Time `json:"createdAt"` // use time.Time to parse ISO 8601 format
}

type PurchaseOrder struct {
	Merchant       Merchant       `json:"-"`
	MerchantShow   interface{}    `json:"merchant"`
	Items          []MerchantItem `json:"items"`
	IsMerchantShow bool           `json:"-"`
}

func (po *PurchaseOrder) SetMerchantShow() {
	if po.IsMerchantShow {
		po.MerchantShow = po.Merchant
	} else {
		po.MerchantShow = nil
	}
}

type OrderDataResponse struct {
	OrderID string          `json:"orderId"`
	Orders  []PurchaseOrder `json:"orders"`
}

type MerchantRequestParams struct {
	MerchantID       string         `query:"merchantId"`
	Limit            int            `query:"limit"`
	Offset           int            `query:"offset"`
	Name             string         `query:"name"`
	MerchantCategory string         `query:"merchantCategory"`
	UserID           string         `query:"-"`
	UserLocation     model.Location `query:"-"`
}

type GetNearbyMerchantResponse struct {
	Distance       float64        `json:"-"`
	Merchant       Merchant       `json:"-"`
	MerchantShow   interface{}    `json:"merchant"`
	Items          []MerchantItem `json:"items"`
	IsMerchantShow bool           `json:"-"`
}

func (po *GetNearbyMerchantResponse) SetDistance(distance float64) {
	po.Distance = distance
}

func (po *GetNearbyMerchantResponse) SetMerchantShow() {
	if po.IsMerchantShow {
		po.MerchantShow = po.Merchant
	} else {
		po.MerchantShow = nil
	}
}
