package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	merchantModel "ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/dto"
	"ps-beli-mang/internal/purchase/model"
	"ps-beli-mang/pkg/errs"
	"strings"
	"time"
)

type orderRepositoryImpl struct {
	db *sqlx.DB
}

func NewOrderRepositoryImpl(db *sqlx.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

const queryCheckMerchantItem = `
	SELECT
		mi.id, mi.merchant_id,mi.price, m.loc_lat, m.loc_long
		FROM merchant_items mi
		JOIN merchants m
		ON mi.merchant_id = m.id
		WHERE mi.merchant_id = ANY($1)
		AND	mi.id = ANY($2)
`

func (o orderRepositoryImpl) GetMerchantItems(ctx context.Context, args []interface{}) ([]merchantModel.MerchantItem, error) {
	result := make([]merchantModel.MerchantItem, 0)
	rows, err := o.db.QueryContext(ctx, queryCheckMerchantItem, args...)
	if err != nil {
		return result, errs.NewErrInternalServerErrors("Error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item merchantModel.MerchantItem
		var merchant merchantModel.Merchant
		if err := rows.Scan(&item.ID, &item.MerchantID, &item.Price, &merchant.LocLat, &merchant.LocLong); err != nil {
			return result, errs.NewErrInternalServerErrors("Error scanning row: %v", err)
		}

		item.SetMerchant(merchant)
		result = append(result, item)
	}

	return result, nil
}

const (
	insertOrderQuery = `
		INSERT INTO orders (id, user_id, total_price, delivery_time, is_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	insertOrderItemQuery = `
		INSERT INTO order_items (id, user_id, order_id, merchant_id, merchant_item_id, quantity, price, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
)

func (o orderRepositoryImpl) SaveOrder(ctx context.Context, order model.Order) error {
	// Start a transaction
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errs.NewErrInternalServerErrors("error starting transaction: %w", err)
	}

	// Insert the order
	_, err = tx.ExecContext(ctx, insertOrderQuery, order.ID, order.UserID, order.TotalPrice, order.DeliveryTime, order.IsOrder, time.Now())
	if err != nil {
		tx.Rollback()
		return errs.NewErrInternalServerErrors("error inserting order: %w", err)
	}

	// Insert the order items
	for _, item := range order.OrderItems {
		_, err = tx.ExecContext(ctx, insertOrderItemQuery, item.ID, item.UserID, item.OrderID, item.MerchantID, item.MerchantItemID, item.Quantity, item.Price, time.Now())
		if err != nil {
			err := tx.Rollback()
			return errs.NewErrInternalServerErrors("error inserting order item: %w", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return errs.NewErrInternalServerErrors("error committing transaction: %w", err)
	}

	return nil
}

const updateUpdateIsOrderTrueQuery = `
		WITH order_to_update AS (
			SELECT id
			FROM orders
			WHERE id = $1
		)
		UPDATE orders
		SET is_order = TRUE
		WHERE id = (SELECT id FROM order_to_update);
	`

func (o orderRepositoryImpl) UpdateOrderSetIsOrderTrue(ctx context.Context, orderID string) error {
	result, err := o.db.ExecContext(ctx, updateUpdateIsOrderTrueQuery, orderID)
	if err != nil {
		return errs.NewErrInternalServerErrors("error updating order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errs.NewErrInternalServerErrors("error checking rows affected: %w", err)
	}

	// If no rows were affected, the order ID does not exist
	if rowsAffected == 0 {
		return errs.NewDefaultErrDataNotFound("calculatedEstimateId is not found")
	}

	// If we reach here, the order was updated successfully
	return nil
}

// Function to build order history query
func buildOrderHistoryQuery(params dto.OrderDataRequestParams) string {
	var filters []string

	// Add conditions based on the parameters
	if params.MerchantID != "" {
		filters = append(filters, fmt.Sprintf("oi.merchant_id = '%s'", params.MerchantID))
	}
	if params.Name != "" {
		filters = append(filters, fmt.Sprintf("(LOWER(m.name) LIKE LOWER('%%%s%%') OR LOWER(mi.name) LIKE LOWER('%%%s%%'))", params.Name, params.Name))
	}
	if params.MerchantCategory != "" {
		filters = append(filters, fmt.Sprintf("m.merchant_category = '%s'", params.MerchantCategory))
	}

	filters = append(filters, "uo.is_order = true")

	// Construct query using CTE
	query := fmt.Sprintf(`
		WITH user_orders AS (
			SELECT o.id AS order_id, o.user_id, o.total_price, o.delivery_time, o.is_order, o.created_at AS order_created_at
			FROM orders o
			WHERE o.user_id = '%s'
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
			mi.category AS merchant_item_category
		FROM 
			user_orders uo
		JOIN 
			order_items oi ON uo.order_id = oi.order_id
		JOIN 
			merchants m ON oi.merchant_id = m.id
		JOIN 
			merchant_items mi ON oi.merchant_item_id = mi.id`, params.UserID)

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY uo.order_created_at DESC"

	limit := 5
	if params.Limit > 0 {
		limit = params.Limit
	}
	offset := 0
	if params.Offset > 0 {
		offset = params.Offset
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	return query
}

func (o orderRepositoryImpl) GetOrdersByUser(ctx context.Context, params dto.OrderDataRequestParams) ([]dto.OrderDataResponse, error) {
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
				Items:          []dto.PurchaseItem{},
				IsMerchantShow: isMerchantShow,
			}

			newPurchaseOrder.SetMerchantShow()

			orderData.Orders = append(orderData.Orders, newPurchaseOrder)
			existingPurchaseOrderIndex = len(orderData.Orders) - 1
		}

		purchaseItem := dto.PurchaseItem{
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

	var results []dto.OrderDataResponse
	for _, orderData := range orderMap {
		results = append(results, *orderData)
	}

	return results, nil
}

func matchesName(name, filter string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(filter))
}
