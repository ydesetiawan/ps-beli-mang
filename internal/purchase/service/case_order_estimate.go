package service

import (
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/pkg/errs"
)

func (o orderService) OrderEstimate(ctx context.Context, request *dto.OrderEstimateRequest) (dto.OrderEstimateResponse, error) {

	merchantIds := make(map[string]struct{})
	itemIds := make(map[string]struct{})
	for _, order := range request.Orders {
		merchantIds[order.MerchantID] = struct{}{}
		for _, item := range order.Items {
			itemIds[item.ItemID] = struct{}{}
		}
	}

	//isStartingPoint not null | there's should be one isStartingPoint == true in orders array
	//if none are true, or true > 1 items, it's not valid
	err := validateOrderData(request)
	if err != nil {
		return dto.OrderEstimateResponse{}, err
	}

	merchantItems, err := o.orderRepository.GetMerchantItem(ctx, buildParams(merchantIds, itemIds))
	if err != nil {
		return dto.OrderEstimateResponse{}, err
	}

	// 404 if merchantId / itemId is not found
	err = IsValidMerchantData(merchantIds, itemIds, merchantItems)
	if err != nil {
		return dto.OrderEstimateResponse{}, err
	}

	return dto.OrderEstimateResponse{}, err

}

func buildParams(merchantIds map[string]struct{}, itemIds map[string]struct{}) []interface{} {
	var args []interface{}

	merchantIdList := keys(merchantIds)
	itemIdList := keys(itemIds)
	args = append(args, pq.Array(merchantIdList), pq.Array(itemIdList))
	return args
}

func validateOrderData(request *dto.OrderEstimateRequest) error {
	// Validate there's exactly one isStartingPoint == true
	startingPoints := 0
	for _, order := range request.Orders {
		if order.IsStartingPoint {
			startingPoints++
		}
	}

	if startingPoints != 1 {
		return errs.NewErrBadRequest("there must be exactly one order with isStartingPoint == true")
	}

	return nil
}

func keys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func IsValidMerchantData(merchantIds, itemIds map[string]struct{}, merchantItems []model.MerchantItem) error {
	foundMerchants := make(map[string]struct{}, len(merchantIds))
	foundItems := make(map[string]struct{}, len(itemIds))
	for _, item := range merchantItems {
		foundItems[item.ID] = struct{}{}
		foundMerchants[item.MerchantID] = struct{}{}
	}
	for merchantId := range merchantIds {
		if _, found := foundItems[merchantId]; !found {
			return errs.NewErrDataNotFound("Merchant ID %s not found in the database", merchantId, errs.ErrorData{})
		}
	}
	for itemId := range itemIds {
		if _, found := foundItems[itemId]; !found {
			return errs.NewErrDataNotFound("Item ID %s not found in the database", itemId, errs.ErrorData{})
		}
	}
	return nil
}
