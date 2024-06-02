package repository

import (
	"context"
	"database/sql"
	"fmt"
	"ps-beli-mang/internal/merchant/dto"
	"ps-beli-mang/internal/merchant/model"
	"ps-beli-mang/pkg/errs"
	"ps-beli-mang/pkg/helper"
	"strings"

	"github.com/jmoiron/sqlx"
)

type merchantRepositoryImpl struct {
	db *sqlx.DB
}

func NewMerchantRepositoryImpl(db *sqlx.DB) MerchantRepository {
	return &merchantRepositoryImpl{db: db}
}

const createMerchantQuery = `INSERT INTO merchants (id, name, merchant_category, image_url, loc_lat, loc_long) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

func (r *merchantRepositoryImpl) CreateMerchant(ctx context.Context, req *dto.MerchantDto) (id string, err error) {
	err = r.db.QueryRowContext(ctx, createMerchantQuery, helper.GenerateULID(), req.Name, req.MerchantCategory, req.ImageUrl, req.Location.Lat, req.Location.Long).Scan(&id)
	return id, err
}

const (
	merchantCteQuery = `WITH MerchantCTE AS (
		SELECT id, name, merchant_category, image_url, loc_lat, loc_long, created_at
		FROM merchants`
	merchantCteTotalQuery = `
		)
		SELECT COUNT(*) OVER() as total, id, name, merchant_category, image_url, loc_lat, loc_long, created_at
		FROM MerchantCTE`
)

func (r *merchantRepositoryImpl) GetMerchants(ctx context.Context, req *dto.MerchantQuery) (merchants []model.Merchant, total int, err error) {
	query := merchantCteQuery
	args := []interface{}{}
	var conditions []string

	if req.ID != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", len(args)+1))
		args = append(args, req.ID)
	}
	if req.Name != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE $%d", len(args)+1))
		args = append(args, "%"+strings.ToLower(req.Name)+"%")
	}
	if req.MerchantCategory != "" {
		conditions = append(conditions, fmt.Sprintf("merchant_category = $%d", len(args)+1))
		args = append(args, req.MerchantCategory)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += merchantCteTotalQuery

	if req.SortCreatedAt != "" {
		if req.SortCreatedAt == model.SortTypeAsc || req.SortCreatedAt == model.SortTypeDesc {
			query += fmt.Sprintf(` ORDER BY created_at %s`, req.SortCreatedAt)
		}
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	args = append(args, req.Limit, req.Offset)

	type MerchantWithTotal struct {
		model.Merchant
		Total int `db:"total"`
	}

	var results []MerchantWithTotal
	err = r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, 0, err
	}

	if len(results) > 0 {
		total = results[0].Total
		for _, result := range results {
			merchants = append(merchants, result.Merchant)
		}
	}

	return merchants, total, err
}

const createMerchantItemQuery = `
    WITH MerchantCheck AS (
        SELECT id
        FROM merchants
        WHERE id = $6
    )
    INSERT INTO merchant_items (id, name, category, image_url, price, merchant_id)
    SELECT $1, $2, $3, $4, $5, id
    FROM MerchantCheck
    RETURNING id;
    `

func (r *merchantRepositoryImpl) CreateMerchantItem(ctx context.Context, merchantId string, req *dto.MerchantItemDto) (id string, err error) {
	err = r.db.QueryRowContext(ctx, createMerchantItemQuery, helper.GenerateULID(), req.Name, req.ProductCategory, req.ImageUrl, req.Price, merchantId).Scan(&id)
	if err != nil {
		if sql.ErrNoRows == err {
			return "", errs.NewErrDataNotFound("No matching merchant found or no row inserted.", req.MerchantId, errs.ErrorData{})
		} else {
			return "", errs.NewErrInternalServerErrors("error insert merchant item", req)
		}
	}
	return id, nil

}

const (
	merchantItemCteQuery = `WITH MerchantItemCTE AS (
		SELECT id, merchant_id, name, category, image_url, price, created_at
		FROM merchant_items`
	merchantItemCteTotalQuery = `
		)
		SELECT COUNT(*) OVER() as total, id, merchant_id, name, category, image_url, price, created_at
		FROM MerchantItemCTE`
)

func (r *merchantRepositoryImpl) GetMerchantItems(ctx context.Context, merchantId string, req *dto.MerchantItemQuery) (merchantItems []model.MerchantItem, total int, err error) {
	query := merchantItemCteQuery
	args := []interface{}{merchantId}
	conditions := []string{fmt.Sprintf("merchant_id = $%d", len(args))}

	if req.ID != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", len(args)+1))
		args = append(args, req.ID)
	}
	if req.Name != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE $%d", len(args)+1))
		args = append(args, "%"+strings.ToLower(req.Name)+"%")
	}
	if req.ProductCategory != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(args)+1))
		args = append(args, req.ProductCategory)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += merchantItemCteTotalQuery

	if req.SortCreatedAt != "" {
		if req.SortCreatedAt == model.SortTypeAsc || req.SortCreatedAt == model.SortTypeDesc {
			query += fmt.Sprintf(` ORDER BY created_at %s`, req.SortCreatedAt)
		}
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	args = append(args, req.Limit, req.Offset)

	type MerchantItemWithTotal struct {
		model.MerchantItem
		Total int `db:"total"`
	}

	var results []MerchantItemWithTotal
	err = r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, 0, err
	}

	if len(results) > 0 {
		total = results[0].Total
		for _, result := range results {
			merchantItems = append(merchantItems, result.MerchantItem)
		}
	}

	return merchantItems, total, err
}
