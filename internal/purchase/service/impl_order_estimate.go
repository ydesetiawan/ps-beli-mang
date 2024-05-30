package service

import (
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"math"
	merchantModel "ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/model"
	"ps-beli-mang/pkg/errs"
	"ps-beli-mang/pkg/helper"
)

func (o orderService) OrderEstimate(ctx context.Context, request dto.OrderEstimateRequest) (dto.OrderEstimateResponse, error) {
	orderId := helper.GenerateULID()
	var response dto.OrderEstimateResponse
	//Validation
	preOrder, err := validationOrderEstimate(ctx, request, o)
	if err != nil {
		return response, err
	}

	//Build Order
	order, err := buildOrder(orderId, preOrder)
	if err != nil {
		return response, err
	}

	//Save Order
	err = o.orderRepository.SaveOrder(ctx, order)
	if err != nil {
		return response, err
	}

	return dto.OrderEstimateResponse{
		TotalPrice:                     order.TotalPrice,
		EstimatedDeliveryTimeInMinutes: order.DeliveryTime,
		CalculatedEstimateId:           order.ID,
	}, err

}

func buildOrder(orderId string, preOrder dto.PreOrder) (model.Order, error) {
	totalPrice := 0.0
	mapMerchantLocation := make(map[string]merchantModel.Location)
	orderItems := make([]model.OrderItem, 0)
	for _, item := range preOrder.MerchantItems {
		totalPrice += item.Price * float64(preOrder.ItemQtyIds[item.ID])
		orderItem := model.OrderItem{
			ID:             helper.GenerateULID(),
			UserID:         preOrder.UserID,
			OrderID:        orderId,
			MerchantID:     item.MerchantID,
			MerchantItemID: item.ID,
			Price:          item.Price,
			Quantity:       preOrder.ItemQtyIds[item.ID],
		}
		orderItems = append(orderItems, orderItem)
		mapMerchantLocation[item.MerchantID] = item.Merchant().Location()
	}

	merchantLocation := getMerchantLocation(mapMerchantLocation, preOrder)

	//EstimateDeliveryTimeTSP
	estimateDeliveryTIme, err := EstimateDeliveryTimeTSP(merchantLocation, preOrder.UserLocation)
	if err != nil {
		return model.Order{}, err
	}

	order := model.Order{
		ID:           orderId,
		UserID:       preOrder.UserID,
		TotalPrice:   totalPrice,
		DeliveryTime: int(math.Round(estimateDeliveryTIme.Minutes())),
		IsOrder:      false,
		OrderItems:   orderItems,
	}
	return order, nil
}

func getMerchantLocation(mapMerchantLocation map[string]merchantModel.Location, orderData dto.PreOrder) []merchantModel.Location {
	merchantLocation := append([]merchantModel.Location{}, mapMerchantLocation[orderData.MerchantStartingPointId])
	for key, item := range mapMerchantLocation {
		if key != orderData.MerchantStartingPointId {
			merchantLocation = append(merchantLocation, item)
		}
	}
	return merchantLocation
}

func validationOrderEstimate(ctx context.Context, request dto.OrderEstimateRequest, o orderService) (dto.PreOrder, error) {
	var preOrder dto.PreOrder
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
		return preOrder, err
	}

	//Get Merchant Item in merchant id and item id
	merchantItems, err := o.orderRepository.GetMerchantItems(ctx, buildParams(merchantIds, itemQtyIds))
	if err != nil {
		return preOrder, err
	}

	// 404 if merchantId / itemId is not found
	err = isValidMerchantData(merchantIds, itemQtyIds, merchantItems)
	if err != nil {
		return preOrder, err
	}

	preOrder.ItemQtyIds = itemQtyIds
	preOrder.MerchantItems = merchantItems
	preOrder.MerchantStartingPointId = merchantStartingPointId
	preOrder.UserLocation = request.UserLocation
	preOrder.UserID = request.UserID
	return preOrder, nil
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
func isValidMerchantData(merchantIds map[string]struct{}, itemIds map[string]int, merchantItems []merchantModel.MerchantItem) error {
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
