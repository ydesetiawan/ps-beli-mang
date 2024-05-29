package model

import (
	"database/sql"
	"time"
)

type MerchantItem struct {
	ID         string         `db:"id"`
	MerchantID string         `db:"merchant_id"`
	Long       float64        `db:"-"`
	Lat        float64        `db:"-"`
	Name       string         `db:"name"`
	Category   sql.NullString `db:"category"`
	ImageURL   sql.NullString `db:"image_url"`
	Quantity   int            `db:"-"`
	Price      float64        `db:"price"`
	CreatedAt  time.Time      `db:"created_at"`
}
