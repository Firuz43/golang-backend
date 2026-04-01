package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type OrderHandler struct {
	DB *sqlx.DB
}

func NewOrderHandler(db *sqlx.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	// 1. Start a Transaction (If one step fails, they all fail)
	tx, err := h.DB.Beginx()
	if err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Safety net

	// 2. Get Cart items and calculate total
	var cartItems []struct {
		ProductID string  `db:"product_id"`
		Quantity  int     `db:"quantity"`
		Price     float64 `db:"price"`
	}

	// Join with products to get the current price
	query := `SELECT c.product_id, c.quantity, p.price FROM cart_items c 
			  JOIN products p ON c.product_id = p.id WHERE c.user_id = $1`

	if err := tx.Select(&cartItems, query, userID); err != nil || len(cartItems) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}

	// 3. Create the Order
	var orderID string
	err = tx.QueryRow("INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id",
		userID, total).Scan(&orderID)
	if err != nil {
		http.Error(w, "Order creation failed", http.StatusInternalServerError)
		return
	}

	// 4. Move items to order_items & Clear Cart
	for _, item := range cartItems {
		tx.Exec("INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase) VALUES ($1, $2, $3, $4)",
			orderID, item.ProductID, item.Quantity, item.Price)
	}

	tx.Exec("DELETE FROM cart_items WHERE user_id = $1", userID)

	// 5. Commit!
	if err := tx.Commit(); err != nil {
		http.Error(w, "Could not complete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"order_id": orderID, "message": "Order placed successfully!"})
}
