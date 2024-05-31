package repository

import (
	"database/sql"
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
