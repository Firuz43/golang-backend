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

	// 2. Inject DB into Handler (Constructor Injection)
	userHandler := handlers.NewUserHandler(db)

	// 3. Routes
	http.HandleFunc("/register", userHandler.RegisterUser)

	http.HandleFunc("/login", userHandler.LoginUser)

	http.HandleFunc("/user", middleware.AuthMiddleware(userHandler.GetUser))

	log.Println("Server is running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
