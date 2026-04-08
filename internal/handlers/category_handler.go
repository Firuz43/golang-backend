package handlers

import (
	"encoding/json"
	"net/http"

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
