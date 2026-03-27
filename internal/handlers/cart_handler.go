package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type CartHandler struct {
	DB *sqlx.DB
}

func NewCartHandler(db *sqlx.DB) *CartHandler {
	return &CartHandler{DB: db}
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// 1. Get UserID from Context (The Bouncer put it there!)
	userID := r.Context().Value("user_id").(string)

	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 2. UPSERT Logic: Insert a new item, OR if it exists, update the quantity
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id) 
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = NOW()
	`

	_, err := h.DB.Exec(query, userID, req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, "Could not add to cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart"})
}
