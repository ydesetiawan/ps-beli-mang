package repository

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/pkg/errs"
)

type orderRepositoryImpl struct {
	db *sqlx.DB
}

func NewOrderRepositoryImpl(db *sqlx.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

const queryCheckMerchantItem = `
SELECT
	mi.id,
	mi.merchant_id,
	m.loc_lat,
	m.loc_long
	FROM
	merchant_items mi
	JOIN
	merchants m
	ON
	mi.merchant_id = m.id
	WHERE
	mi.merchant_id = ANY($1)
	AND
	mi.id = ANY($2)
`

func (o orderRepositoryImpl) GetMerchantItems(ctx context.Context, args []interface{}) ([]model.MerchantItem, error) {
	result := make([]model.MerchantItem, 0)
	rows, err := o.db.QueryContext(ctx, queryCheckMerchantItem, args...)
	if err != nil {
		return result, errs.NewErrInternalServerErrors("Error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.MerchantItem
		var merchant model.Merchant
		if err := rows.Scan(&item.ID, &item.MerchantID, &merchant.LocLat, &merchant.LocLong); err != nil {
			return result, errs.NewErrInternalServerErrors("Error scanning row: %v", err)
		}

		item.SetMerchant(merchant)
		result = append(result, item)
	}

	return result, nil
}
