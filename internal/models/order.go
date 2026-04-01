package models

import "time"

type Order struct {
	ID         string      `json:"id" db:"id"`
	UserID     string      `json:"user_id" db:"user_id"`
	TotalPrice float64     `json:"total_price" db:"total_price"`
	Status     string      `json:"status" db:"status"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	Items      []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID              string  `json:"id" db:"id"`
	OrderID         string  `json:"order_id" db:"order_id"`
	ProductID       string  `json:"product_id" db:"product_id"`
	Quantity        int     `json:"quantity" db:"quantity"`
	PriceAtPurchase float64 `json:"price_at_purchase" db:"price_at_purchase"`
}
