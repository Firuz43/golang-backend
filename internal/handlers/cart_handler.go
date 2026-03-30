package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Firuz43/ecommerce/internal/models"
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

// ######################## G E T C A R T ###########################

// GetCart returns all items in the current user's cart with product details
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	// We use a JOIN to get the Product name and price in one go
	query := `
		SELECT 
			c.id, c.user_id, c.product_id, c.quantity, c.created_at,
			p.name AS "product.name", 
			p.price AS "product.price", 
			p.image_url AS "product.image_url"
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1
	`

	var items []models.CartItem
	err := h.DB.Select(&items, query, userID)
	if err != nil {
		http.Error(w, "Could not fetch cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// RemoveFromCart deletes a specific product from the user's cart
func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	// We can get the product_id from the URL query or a JSON body.//
	// Let's use a URL query for a change: /cart/remove?product_id=UUID
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`

	result, err := h.DB.Exec(query, userID, productID)
	if err != nil {
		http.Error(w, "Delete failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if we actually deleted anything
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Item not found in cart", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart"})
}
