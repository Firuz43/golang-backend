package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Firuz43/ecommerce/internal/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *sqlx.DB //Database connection can be shared across all handler methods
}

// NewUserHandler is our "Constructor"
func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{DB: db}
	//receive the DB connection and store it in the struct for later use
}

type RegisterRequest struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the incoming JSON request
	var req RegisterRequest
	// json.NewDecoder is like Jackson in Java; it reads the request body into our struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 2. Hash the password
	// Bcrypt.DefaultCost is a balance between security and speed
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 3. Save to Database
	// We use 'QueryRow' because we want Postgres to return the ID it generated
	var user models.User
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at` // save email and hashed password, then return the new user's ID, email, and creation time

	err = h.DB.QueryRow(query, req.Email, string(hashedPassword)).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		// In production, check if error is "unique_violation" (email already exists)
		http.Error(w, "Could not create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Respond to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Send 201 Created status
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// A simple query to get the first user from the database//
	// Coming from Java? Notice how we don't need a heavy ORM here!!
	err := h.DB.Get(&user, "SELECT id, email, created_at FROM users LIMIT 1")

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
