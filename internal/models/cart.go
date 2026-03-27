package models

import "time"

type CartItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json: "created_at" db:"created_at"`

	Product *Product `json:"product, omitempty"`
}
