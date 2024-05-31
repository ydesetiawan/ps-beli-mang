package model

import (
	"slices"
	"time"
)

type ItemCategory string

const (
	Beverage   ItemCategory = "Beverage"
	Food       ItemCategory = "Food"
	Snack      ItemCategory = "Snack"
	Condiments ItemCategory = "Condiments"
	Additions  ItemCategory = "Additions"
)

var itemCategories = []ItemCategory{Beverage, Food, Snack, Condiments, Additions}

func (pc ItemCategory) String() string {
	return string(pc)
}

func (pc ItemCategory) Valid() bool {
	return slices.Index(itemCategories, ItemCategory(pc)) != -1
}

func ValidMerchantItemCategory(category string) bool {
	var categories = []ItemCategory{
		Beverage,
		Food,
		Snack,
		Condiments,
		Additions,
	}

	return slices.IndexFunc(categories, func(ic ItemCategory) bool {
		return ic.String() == category
	}) != -1
}

type MerchantItem struct {
	ID         string
	Name       string
	Category   ItemCategory
	ImageUrl   string
	Price      float64
	CreatedAt  time.Time
	CreatedBy  string
	MerchantID string
	merchant   Merchant
}

func (m *MerchantItem) SetMerchant(merchant Merchant) {
	m.merchant = merchant
}

func (m MerchantItem) Merchant() Merchant {
	return m.merchant
}
