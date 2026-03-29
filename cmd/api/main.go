package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Firuz43/ecommerce/internal/database"
	"github.com/Firuz43/ecommerce/internal/handlers"
	"github.com/Firuz43/ecommerce/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {

	// 1. Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Now you can access variables using os.Getenv
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	log.Printf("Server starting on port %s...", port)

	// 1. Initialize Database (This also runs migrations)
	db, err := database.NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	//Initialize Handlers
	// 2. Inject DB into Handler (Constructor Injection)
	userHandler := handlers.NewUserHandler(db)
	productHandler := handlers.NewProductHandler(db)
	cartHandler := handlers.NewCartHandler(db)

	// 3. Routes
	//Auth routes
	http.HandleFunc("/register", userHandler.RegisterUser)
	http.HandleFunc("/login", userHandler.LoginUser)
	http.HandleFunc("/user", middleware.AuthMiddleware(userHandler.GetUser))

	// Product routes
	http.HandleFunc("/products", productHandler.GetProducts)
	http.HandleFunc("/products/add", middleware.AuthMiddleware(productHandler.CreateProduct))

	// Cart routes
	// POST /cart - Adds an item to the user's specific cart
	// http.HandleFunc("/cart", middleware.AuthMiddleware(cartHandler.AddToCart))//
	http.HandleFunc("/cart", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cartHandler.AddToCart(w, r)
		} else if r.Method == http.MethodGet {
			cartHandler.GetCart(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	log.Println("Server is running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
