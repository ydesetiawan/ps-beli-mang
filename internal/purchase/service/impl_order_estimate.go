package service

import (
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"math"
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/pkg/errs"
	"ps-beli-mang/pkg/helper"
)

func (o orderService) OrderEstimate(ctx context.Context, request dto.OrderEstimateRequest) (dto.OrderEstimateResponse, error) {
	id := helper.GenerateULID()
	var result dto.OrderEstimateResponse
	//Validation
	orderData, err := validationOrderEstimate(ctx, request, o)
	if err != nil {
		return result, err
	}

	//Calculate total Price
	totalPrice := 0
	mapMerchantLocation := make(map[string]model.Location)
	for _, item := range orderData.MerchantItems {
		totalPrice += item.Price * orderData.ItemQtyIds[item.ID]
		mapMerchantLocation[item.MerchantID] = item.Merchant().Location()
	}

	//EstimateDeliveryTimeTSP
	merchantLocation := getMerchantLocation(mapMerchantLocation, orderData)
	estimateDeliveryTIme, err := EstimateDeliveryTimeTSP(merchantLocation, request.UserLocation)
	if err != nil {
		return result, err
	}

	//TODO save to order and order detail

	return dto.OrderEstimateResponse{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: int(math.Round(estimateDeliveryTIme.Minutes())),
		CalculatedEstimateId:           id,
	}, err

}

func getMerchantLocation(mapMerchantLocation map[string]model.Location, orderData dto.OrderEstimateProcess) []model.Location {
	merchantLocation := append([]model.Location{}, mapMerchantLocation[orderData.MerchantStartingPointId])
	for key, item := range mapMerchantLocation {
		if key != orderData.MerchantStartingPointId {
			merchantLocation = append(merchantLocation, item)
		}
	}
	return merchantLocation
}

func validationOrderEstimate(ctx context.Context, request dto.OrderEstimateRequest, o orderService) (dto.OrderEstimateProcess, error) {
	var orderData dto.OrderEstimateProcess
	merchantIds := make(map[string]struct{})
	itemQtyIds := make(map[string]int)
	for _, order := range request.Orders {
		merchantIds[order.MerchantID] = struct{}{}
		for _, item := range order.Items {
			itemQtyIds[item.ItemID] = item.Quantity
		}
	}

	//isStartingPoint not null | there's should be one isStartingPoint == true in orders array
	//if none are true, or true > 1 items, it's not valid
	merchantStartingPointId, err := validateStartingPoints(request)
	if err != nil {
		return orderData, err
	}

	//Get Merchant Item in merchant id and item id
	merchantItems, err := o.orderRepository.GetMerchantItems(ctx, buildParams(merchantIds, itemQtyIds))
	if err != nil {
		return orderData, err
	}

	// 404 if merchantId / itemId is not found
	err = isValidMerchantData(merchantIds, itemQtyIds, merchantItems)
	if err != nil {
		return orderData, err
	}

	orderData.ItemQtyIds = itemQtyIds
	orderData.MerchantItems = merchantItems
	orderData.MerchantStartingPointId = merchantStartingPointId
	return orderData, nil
}

func buildParams(merchantIds map[string]struct{}, itemIds map[string]int) []interface{} {
	var args []interface{}

	merchantIdList := keys(merchantIds)
	itemIdList := intKeys(itemIds)
	args = append(args, pq.Array(merchantIdList), pq.Array(itemIdList))
	return args
}

// Return merchantStartingPointId, error
func validateStartingPoints(request dto.OrderEstimateRequest) (string, error) {
	// Validate there's exactly one isStartingPoint == true
	startingPoints := 0
	var merchantStartingPointId string
	for _, order := range request.Orders {
		if order.IsStartingPoint {
			merchantStartingPointId = order.MerchantID
			startingPoints++
		}
	}

	if startingPoints != 1 {
		return merchantStartingPointId, errs.NewErrBadRequest("there must be exactly one order with isStartingPoint == true")
	}

	return merchantStartingPointId, nil
}

func keys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func intKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Validate Merchant ID and Item ID
func isValidMerchantData(merchantIds map[string]struct{}, itemIds map[string]int, merchantItems []model.MerchantItem) error {
	foundMerchants := make(map[string]struct{}, len(merchantIds))
	foundItems := make(map[string]struct{}, len(itemIds))
	for _, item := range merchantItems {
		foundItems[item.ID] = struct{}{}
		foundMerchants[item.MerchantID] = struct{}{}
	}
	for merchantId := range merchantIds {
		if _, found := foundMerchants[merchantId]; !found {
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
