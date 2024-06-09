package repository

import (
	"fmt"
	"golang.org/x/net/context"
	merchantModel "ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"strings"
	"time"
)

// Function to build order history query
func buildOrderHistoryQuery(params dto.MerchantRequestParams) string {
	var filters []string

	// Add conditions based on the parameters
	if params.MerchantID != "" {
		filters = append(filters, fmt.Sprintf("oi.merchant_id = '%s'", params.MerchantID))
	}
	if params.Name != "" {
		filters = append(filters, fmt.Sprintf("(LOWER(m.name) LIKE LOWER('%%%s%%') AND LOWER(mi.name) LIKE LOWER('%%%s%%'))", params.Name, params.Name))
	}
	if params.MerchantCategory != "" {
		filters = append(filters, fmt.Sprintf("m.merchant_category = '%s'", params.MerchantCategory))
	}

	filters = append(filters, "uo.is_order = true")
	limit := 5
	if params.Limit > 0 {
		limit = params.Limit
	}
	offset := 0
	if params.Offset > 0 {
		offset = params.Offset
	}
	// Construct query using CTE
	query := fmt.Sprintf(`
		WITH user_orders AS (
			SELECT o.id AS order_id, o.user_id, o.total_price, o.delivery_time, o.is_order, o.created_at AS order_created_at
			FROM orders o
			WHERE o.user_id = '%s' 
 			LIMIT '%d' OFFSET '%d'
		)
		SELECT 
			uo.order_id,
			uo.user_id,
			uo.total_price,
			uo.delivery_time,
			uo.is_order,
			uo.order_created_at,
			oi.id AS order_item_id,
			oi.merchant_id,
			oi.merchant_item_id,
			oi.quantity,
			oi.price,
			oi.created_at AS order_item_created_at,
			m.name AS merchant_name,
			m.merchant_category,
			m.loc_lat AS latitude,
			m.loc_long AS longitude,
			m.image_url AS merchant_image_url,
			mi.name AS merchant_item_name,
			mi.image_url AS merchant_item_image_url,
			mi.category AS merchant_item_category
		FROM 
			user_orders uo
		JOIN 
			order_items oi ON uo.order_id = oi.order_id
		JOIN 
			merchants m ON oi.merchant_id = m.id
		JOIN 
			merchant_items mi ON oi.merchant_item_id = mi.id`, params.UserID, limit, offset)

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY uo.order_created_at DESC"

	return query
}

func (o orderRepositoryImpl) GetOrdersByUser(ctx context.Context, params dto.MerchantRequestParams) ([]dto.OrderDataResponse, error) {
	query := buildOrderHistoryQuery(params)

	var rawResults []struct {
		OrderID              string    `db:"order_id"`
		UserID               string    `db:"user_id"`
		TotalPrice           float64   `db:"total_price"`
		DeliveryTime         int       `db:"delivery_time"`
		IsOrder              bool      `db:"is_order"`
		OrderCreatedAt       time.Time `db:"order_created_at"`
		OrderItemID          string    `db:"order_item_id"`
		MerchantID           string    `db:"merchant_id"`
		MerchantName         string    `db:"merchant_name"`
		MerchantCategory     string    `db:"merchant_category"`
		MerchantImageURL     string    `db:"merchant_image_url"`
		Latitude             float64   `db:"latitude"`
		Longitude            float64   `db:"longitude"`
		MerchantCreatedAt    time.Time `db:"merchant_created_at"`
		MerchantItemID       string    `db:"merchant_item_id"`
		MerchantItemName     string    `db:"merchant_item_name"`
		MerchantItemCategory string    `db:"merchant_item_category"`
		MerchantItemImageURL string    `db:"merchant_item_image_url"`
		Price                float64   `db:"price"`
		Quantity             int       `db:"quantity"`
		OrderItemCreatedAt   time.Time `db:"order_item_created_at"`
	}

	err := o.db.SelectContext(ctx, &rawResults, query)
	if err != nil {
		return nil, err
	}

	// Process raw results into the structured response
	orderMap := make(map[string]*dto.OrderDataResponse)
	for _, raw := range rawResults {
		if _, exists := orderMap[raw.OrderID]; !exists {
			orderMap[raw.OrderID] = &dto.OrderDataResponse{
				OrderID: raw.OrderID,
				Orders:  []dto.PurchaseOrder{},
			}
		}

		orderData := orderMap[raw.OrderID]

		var existingPurchaseOrderIndex int
		var existingPurchaseOrderFound bool
		for i := range orderData.Orders {
			if orderData.Orders[i].Merchant.MerchantID == raw.MerchantID {
				existingPurchaseOrderIndex = i
				existingPurchaseOrderFound = true
				break
			}
		}

		if !existingPurchaseOrderFound {
			newMerchant := dto.Merchant{
				MerchantID:       raw.MerchantID,
				Name:             raw.MerchantName,
				MerchantCategory: raw.MerchantCategory,
				ImageURL:         raw.MerchantImageURL,
				Location: merchantModel.Location{
					Lat:  raw.Latitude,
					Long: raw.Longitude,
				},
				CreatedAt: raw.MerchantCreatedAt,
			}

			isMerchantShow := true
			if "" != params.Name && !matchesName(raw.MerchantName, params.Name) {
				isMerchantShow = false
			}

			newPurchaseOrder := dto.PurchaseOrder{
				Merchant:       newMerchant,
				Items:          []dto.MerchantItem{},
				IsMerchantShow: isMerchantShow,
			}

			newPurchaseOrder.SetMerchantShow()

			orderData.Orders = append(orderData.Orders, newPurchaseOrder)
			existingPurchaseOrderIndex = len(orderData.Orders) - 1
		}

		purchaseItem := dto.MerchantItem{
			ItemID:          raw.MerchantItemID,
			Name:            raw.MerchantItemName,
			ProductCategory: raw.MerchantItemCategory,
			Price:           raw.Price,
			Quantity:        raw.Quantity,
			ImageURL:        raw.MerchantItemImageURL,
			CreatedAt:       raw.OrderItemCreatedAt,
		}

		if "" != params.Name && !matchesName(raw.MerchantItemName, params.Name) {
			//Nothing
		} else {
			orderData.Orders[existingPurchaseOrderIndex].Items = append(orderData.Orders[existingPurchaseOrderIndex].Items, purchaseItem)
		}

	}

	results := []dto.OrderDataResponse{}
	for _, orderData := range orderMap {
		results = append(results, *orderData)
	}

	return results, nil
}

func matchesName(name, filter string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(filter))
}
