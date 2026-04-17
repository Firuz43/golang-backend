package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Firuz43/ecommerce/internal/auth"
	"github.com/Firuz43/ecommerce/internal/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Wrap our DB with a struct to securely pass it around to our handler methods
type UserHandler struct {
	DB *sqlx.DB //Database connection can be shared across all handler methods
}

// NewUserHandler is our "Constructor"
func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{DB: db}
	//receive the DB connection and store it in the struct for later use
}

// RegisterRequest defines what data we expect from ..the frontend (Flutter)
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
	// Bcrypt.DefaultCost is a balance between security and speed//
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 3. Save to Database
	// We use 'QueryRow' because we want Postgres to return the ID it generated
	var user models.User
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at` // save email and hashed password, then return the new user's ID, email, and creation time//

	err = h.DB.QueryRow(query, req.Email, string(hashedPassword)).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		// In production, check if error is "unique_violation" (email already exists)
		http.Error(w, "Could not create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Respond to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Send 201 Created status//
	json.NewEncoder(w).Encode(user)
}

// ########################################################################## ##########################################################################

// ########################################################################## LOGIN HANDLER ##########################################################################

// LoginRequest captures the email and password from the JSON body
type LoginRequest struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// 1. Setup a struct to capture the incoming JSON (Email/Password)
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 2. Decode the Request Body into our struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Find the user in the database by their email
	var user models.User
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	err := h.DB.Get(&user, query, req.Email)
	if err != nil {
		// Security Tip: Use a generic error so hackers don't know if the email exists
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 4. Compare the plain-text password from the user with the hash from the DB
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 5. SUCCESS! Now generate the "VIP Pass" (JWT)
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// 6. Send the token back to your Flutter app as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"message": "Login successful!",
	})
}

// ########################################################################## GET USER HANDLER ##########################################################################
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 1. Get the UserID from the context (set by the middleware)
	// We use type assertion .(string) because context stores values as interface{}
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Could not get user from context", http.StatusInternalServerError)
		return
	}

	// 2. Fetch the user from the database using that ID
	var user models.User

	// We only select id, email, and created_at because we don't want to send the password hash back to the client
	query := `SELECT id, email, created_at FROM users WHERE id = $1`
	// sqlx's Get method is a convenient way to run a query and scan the result into a struct
	err := h.DB.Get(&user, query, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 3. Send the user data back (without the password hash!)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
