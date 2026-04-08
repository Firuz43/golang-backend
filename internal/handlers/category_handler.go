package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Firuz43/ecommerce/internal/models"
	"github.com/jmoiron/sqlx"
)

type CategoryHandler struct {
	DB *sqlx.DB
}

func NewCategoryHandler(db *sqlx.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	var id string
	err := h.DB.QueryRow(query, req.Name, req.Description).Scan(&id)
	if err != nil {
		http.Error(w, "Could not create category", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": id, "message": "Category created!"})
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	var cats []interface{} // You can use your model here
	err := h.DB.Select(&cats, "SELECT * FROM categories")
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cats)
}

// Updated GetProducts to allow filtering by Category ID
func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	var products []models.Product
	var err error

	if categoryID != "" {
		// Filtered view
		query := `SELECT * FROM products WHERE category_id = $1 ORDER BY created_at DESC`
		err = h.DB.Select(&products, query, categoryID)
	} else {
		// General view
		query := `SELECT * FROM products ORDER BY created_at DESC`
		err = h.DB.Select(&products, query)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}
