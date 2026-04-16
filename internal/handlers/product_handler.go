package handlers

import (
	"encoding/json"
	"log"
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
// Updated GetProducts to allow filtering by Category ID
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	var products []models.Product
	var err error

	columns := `id, name, description, price, stock, image_url, category_id, created_at`

	if categoryID != "" {
		// Filtered view
		query := `SELECT ` + columns + ` FROM products WHERE category_id = $1 ORDER BY created_at DESC`
		err = h.DB.Select(&products, query, categoryID)
	} else {
		// General view
		query := `SELECT ` + columns + ` FROM products ORDER BY created_at DESC`
		err = h.DB.Select(&products, query)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("GET PRODUCTS ERROR: %v", err)
		return
	}
	json.NewEncoder(w).Encode(products)
}

// ################# C R E A T E  P R O D U C T #################/
// CreateProduct allows adding a new item to the catalog//
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
