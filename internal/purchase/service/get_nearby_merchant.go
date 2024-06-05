package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/purchase/dto"
	"sort"
	"strings"
)

func (o orderService) GetNearbyMerchants(ctx context.Context, params dto.MerchantRequestParams) ([]dto.GetNearbyMerchantResponse, int, error) {

	merchants, err := o.orderRepository.GetAllMerchants(ctx)
	if err != nil {
		return nil, 0, err
	}

	merchants = filterMerchants(merchants, params)
	merchants = filterMerchantsWithOffsetAndLimitAndSortDistance(merchants, params)

	return merchants, len(merchants), nil
}

func filterMerchants(merchants []dto.GetNearbyMerchantResponse, params dto.MerchantRequestParams) []dto.GetNearbyMerchantResponse {
	var filtered []dto.GetNearbyMerchantResponse

	for _, merchant := range merchants {
		merchant.IsMerchantShow = true
		if params.MerchantID != "" && merchant.Merchant.MerchantID != params.MerchantID {
			continue
		}

		if params.Name != "" {
			nameMatch := strings.Contains(strings.ToLower(merchant.Merchant.Name), strings.ToLower(params.Name))
			itemNameMatch := false
			var newItems = make([]dto.MerchantItem, 0)
			for _, item := range merchant.Items {
				if strings.Contains(strings.ToLower(item.Name), strings.ToLower(params.Name)) {
					itemNameMatch = true
					newItems = append(newItems, item)
				}
			}

			if !nameMatch && !itemNameMatch {
				continue
			}
			merchant.Items = newItems
			if !nameMatch {
				merchant.IsMerchantShow = false
			}
		}

		if params.MerchantCategory != "" && merchant.Merchant.MerchantCategory != params.MerchantCategory {
			continue
		}

		distance := HaversineDistance(params.UserLocation.Lat, params.UserLocation.Long, merchant.Merchant.Location.Lat, merchant.Merchant.Location.Long)
		merchant.SetDistance(distance)
		merchant.SetMerchantShow()
		filtered = append(filtered, merchant)
	}

	return filtered
}

func filterMerchantsWithOffsetAndLimitAndSortDistance(filtered []dto.GetNearbyMerchantResponse, params dto.MerchantRequestParams) []dto.GetNearbyMerchantResponse {

	//Sort Distance
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Distance < filtered[j].Distance
	})

	offset := params.Offset
	if offset > len(filtered) {
		offset = len(filtered)
	}

	limit := params.Limit
	if limit == 0 {
		limit = 5
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end]
}
