package repository

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"time"
)

const getAllMerchantQuery = `
		SELECT 
			m.id AS merchant_id, 
			m.name AS merchant_name, 
			m.merchant_category, 
			m.image_url AS merchant_image_url, 
			m.loc_lat AS merchant_lat, 
			m.loc_long AS merchant_long, 
			m.created_at AS merchant_created_at, 
			mi.id AS item_id, 
			mi.name AS item_name, 
			mi.category AS item_category, 
			mi.image_url AS item_image_url, 
			mi.price AS item_price, 
			mi.created_at AS item_created_at 
		FROM 
			merchants m 
		LEFT JOIN 
			merchant_items mi 
		ON 
			m.id = mi.merchant_id
	`

func (o orderRepositoryImpl) GetNearbyMerchantByUser(ctx context.Context, params dto.MerchantRequestParams) ([]dto.GetNearbyMerchantResponse, error) {
	return o.getAllMerchants(ctx)
}

func (o orderRepositoryImpl) getAllMerchants(ctx context.Context) ([]dto.GetNearbyMerchantResponse, error) {

	var rawResults []struct {
		MerchantID        string    `db:"merchant_id"`
		MerchantName      string    `db:"merchant_name"`
		MerchantCategory  string    `db:"merchant_category"`
		MerchantImageURL  string    `db:"merchant_image_url"`
		Latitude          float64   `db:"merchant_lat"`
		Longitude         float64   `db:"merchant_long"`
		MerchantCreatedAt time.Time `db:"merchant_created_at"`
		ItemID            string    `db:"item_id"`
		ItemName          string    `db:"item_name"`
		ItemCategory      string    `db:"item_category"`
		ItemImageURL      string    `db:"item_image_url"`
		Price             float64   `db:"item_price"`
		ItemCreatedAt     time.Time `db:"item_created_at"`
	}

	err := o.db.SelectContext(ctx, &rawResults, getAllMerchantQuery)
	if err != nil {
		return nil, err
	}

	merchantsMap := make(map[string]*dto.GetNearbyMerchantResponse)
	for _, rawResult := range rawResults {
		merchantID := rawResult.MerchantID

		// Check if the merchant already exists in the map
		if _, exists := merchantsMap[merchantID]; !exists {
			// Create a new merchant entry
			merchant := dto.Merchant{
				MerchantID:       rawResult.MerchantID,
				Name:             rawResult.MerchantName,
				MerchantCategory: rawResult.MerchantCategory,
				ImageURL:         rawResult.MerchantImageURL,
				Location: model.Location{
					Lat:  rawResult.Latitude,
					Long: rawResult.Longitude,
				},
				CreatedAt: rawResult.MerchantCreatedAt,
			}

			// Create a new response object
			resp := dto.GetNearbyMerchantResponse{
				Merchant: merchant,
				Items:    []dto.MerchantItem{},
			}

			// Add to the map
			merchantsMap[merchantID] = &resp
		}

		// Append items to the merchant
		if rawResult.ItemID != "" {
			item := dto.MerchantItem{
				ItemID:          rawResult.ItemID,
				Name:            rawResult.ItemName,
				ProductCategory: rawResult.ItemCategory,
				Price:           rawResult.Price,
				ImageURL:        rawResult.ItemImageURL,
				CreatedAt:       rawResult.ItemCreatedAt,
			}

			// Append the item to the corresponding merchant's items slice
			merchantsMap[merchantID].Items = append(merchantsMap[merchantID].Items, item)
		}
	}

	var response []dto.GetNearbyMerchantResponse
	for _, merchant := range merchantsMap {
		response = append(response, *merchant)
	}

	return response, nil
}
