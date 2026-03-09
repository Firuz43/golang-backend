package main

import (
	"log"
	"net/http"

	"github.com/Firuz43/ecommerce/internal/database"
	"github.com/Firuz43/ecommerce/internal/handlers"
)

func main() {

	// 1. Initialize Database (This also runs migrations)
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// 2. Inject DB into Handler (Constructor Injection)
	userHandler := handlers.NewUserHandler(db)

	// 3. Routes
	http.HandleFunc("/register", userHandler.RegisterUser)

	http.HandleFunc("/user", userHandler.GetUser)

	log.Println("Server is running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
