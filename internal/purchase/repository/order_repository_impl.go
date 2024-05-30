package repository

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	merchantModel "ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/internal/purchase/model"
	"ps-beli-mang/pkg/errs"
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
		INSERT INTO orders (id, total_price, delivery_time, is_order, user_loc_lat, user_loc_long, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	insertOrderItemQuery = `
		INSERT INTO order_items (id, order_id, merchant_id, is_starting_point, merchant_item_id, quantity, price, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
)

func (o orderRepositoryImpl) SaveOrder(ctx context.Context, order model.Order) error {
	// Insert the order
	_, err := o.db.ExecContext(ctx, insertOrderQuery, order.ID, order.TotalPrice, order.DeliveryTime, order.IsOrder, order.UserLocLat, order.UserLocLong, time.Now())
	if err != nil {
		return errs.NewErrInternalServerErrors("Error SaveOrder: %v", err)
	}
	// Insert the order items
	for _, item := range order.OrderItems {
		_, err = o.db.ExecContext(ctx, insertOrderItemQuery, item.ID, item.OrderID, item.MerchantID, item.IsStartingPoint, item.MerchantItemID, item.Quantity, item.Price, item.Amount, time.Now())
		if err != nil {
			return errs.NewErrInternalServerErrors("Error SaveOrderItem: %v", err)
		}
	}

	return nil
}
