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

// ################ G E T  A L L  P R O D U C T S #################
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

// ################# C R E A T E  P R O D U C T #################
// CreateProduct allows adding a new item to the catalog
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Define the input structure
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		ImageURL    string  `json:"image_url"`
	}

	// 2. Decode the incoming JSON''
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 3. Insert into Database
	// We let Postgres handle the ID (UUID) and the Timestamps automatically
	query := `INSERT INTO products (name, description, price, stock, image_url) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := h.DB.QueryRow(query, req.Name, req.Description, req.Price, req.Stock, req.ImageURL).Scan(&id)

	if err != nil {
		http.Error(w, "Failed to save product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Product created successfully!",
	})
}
