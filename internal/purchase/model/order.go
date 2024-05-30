package model

import "time"

type Order struct {
	ID           string      `db:"id"`
	TotalPrice   float64     `db:"total_price"`
	DeliveryTime int         `db:"delivery_time"`
	IsOrder      bool        `db:"is_order"`
	UserLocLat   float64     `db:"user_loc_lat"`
	UserLocLong  float64     `db:"user_loc_long"`
	CreatedAt    time.Time   `db:"created_at"`
	OrderItems   []OrderItem `db:"-"`
}

type OrderItem struct {
	ID              string    `db:"id"`
	OrderID         string    `db:"order_id"`
	MerchantID      string    `db:"merchant_id"`
	IsStartingPoint bool      `db:"is_starting_point"`
	MerchantItemID  string    `db:"merchant_item_id"`
	Quantity        int       `db:"quantity"`
	Price           float64   `db:"price"`
	Amount          float64   `db:"amount"`
	CreatedAt       time.Time `db:"created_at"`
}
