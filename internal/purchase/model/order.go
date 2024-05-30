package model

import "time"

type Order struct {
	ID           string      `db:"id"`
	UserID       string      `db:"user_id"`
	TotalPrice   float64     `db:"total_price"`
	DeliveryTime int         `db:"delivery_time"`
	IsOrder      bool        `db:"is_order"`
	CreatedAt    time.Time   `db:"created_at"`
	OrderItems   []OrderItem `db:"-"`
}

type OrderItem struct {
	ID             string    `db:"id"`
	UserID         string    `db:"user_id"`
	OrderID        string    `db:"order_id"`
	MerchantID     string    `db:"merchant_id"`
	MerchantItemID string    `db:"merchant_item_id"`
	Quantity       int       `db:"quantity"`
	Price          float64   `db:"price"`
	CreatedAt      time.Time `db:"created_at"`
}
