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

var merchantCategories = []MerchantCategory{
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
	ID        string           `db:"id"`
	Name      string           `db:"name"`
	Category  MerchantCategory `db:"merchant_category"`
	ImageUrl  string           `db:"image_url"`
	LocLat    float64          `db:"loc_lat"`
	LocLong   float64          `db:"loc_long"`
	CreatedAt time.Time        `db:"created_at"`
}

type SortType string

const CACHE_KEY_ALL_MERCHANTS = "all_merchants"

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
