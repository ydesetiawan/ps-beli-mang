package model

import (
	"slices"
	"time"
)

type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

var merchantCategories []MerchantCategory = []MerchantCategory{
	SmallRestaurant,
	MediumRestaurant,
	LargeRestaurant,
	MerchandiseRestaurant,
	BoothKiosk,
	ConvenienceStore,
}

func (mc MerchantCategory) String() string {
	return string(mc)
}

func (mc MerchantCategory) Valid() bool {
	return slices.Index(merchantCategories, MerchantCategory(mc)) != -1
}

type Merchant struct {
	ID        string
	Name      string
	Category  MerchantCategory
	ImageUrl  string
	LocLat    float64
	LocLong   float64
	CreatedAt time.Time
}

type SortType string

const (
	SortTypeAsc  SortType = "asc"
	SortTypeDesc SortType = "desc"
)

type MerchantFetchFilter struct {
	ID               string           `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory MerchantCategory `json:"merchantCategory"`
	SortCreatedAt    SortType         `json:"createdAt"`
	Limit            int              `json:"limit"`
	Offset           int              `json:"offset"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (mc Merchant) Location() Location {
	return Location{
		Lat:  mc.LocLat,
		Long: mc.LocLong,
	}
}
