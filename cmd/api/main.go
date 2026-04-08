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
	orderHandler := handlers.NewOrderHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)

	// 3. Routes
	//Auth routes
	http.HandleFunc("/register", userHandler.RegisterUser)
	http.HandleFunc("/login", userHandler.LoginUser)
	http.HandleFunc("/user", middleware.AuthMiddleware(userHandler.GetUser))

	// Product routes
	http.HandleFunc("/products", productHandler.GetProducts)
	http.HandleFunc("/products/add", middleware.AuthMiddleware(productHandler.CreateProduct))

	// Cart routes
	// We use a single endpoint for cart operations, and switch based on the HTTP method. This keeps our API clean and RESTful.
	http.HandleFunc("/cart", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cartHandler.GetCart(w, r)
		case http.MethodPost:
			cartHandler.AddToCart(w, r)
		case http.MethodDelete:
			cartHandler.RemoveFromCart(w, r) // New functionality!
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Order routes
	http.HandleFunc("/checkout", middleware.AuthMiddleware(orderHandler.Checkout))
	http.HandleFunc("/orders", middleware.AuthMiddleware(orderHandler.GetOrders))

	// Category routes
	http.HandleFunc("/categories", categoryHandler.GetCategories)
	http.HandleFunc("/categories/add", middleware.AuthMiddleware(categoryHandler.CreateCategory))

	log.Println("Server is running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
