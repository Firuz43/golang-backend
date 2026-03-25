package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Firuz43/ecommerce/internal/models"
	"github.com/jmoiron/sqlx"
)

// ProductHandler will handle all product-related HTTP requests
type ProductHandler struct {
	DB *sqlx.DB
}

// NewProductHandler is our "Constructor" for ProductHandler
func NewProductHandler(db *sqlx.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

// GetProducts returns all products in the database
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	//
	var products []models.Product

	// We want to return products sorted by newest first, so we order by created_at DESC
	query := `SELECT * FROM products ORDER BY created_at DESC`
	// sqlx's Select method will automatically map the rows to our Product struct slice
	err := h.DB.Select(&products, query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
