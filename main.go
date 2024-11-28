package main

import (
	"food-court/handlers"
	"food-court/middleware"
	"food-court/models"
	"food-court/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&models.User{}, &models.MenuItem{}, &models.Order{}, &models.OrderItem{}, &models.Payment{}, &models.Admin{})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	menuHandler := handlers.NewMenuHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	// Initialize router
	r := mux.NewRouter()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},  // Your React app URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},  // Allow all headers
		AllowCredentials: true,
	})

	// Public routes
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Menu routes
	api.HandleFunc("/menu", menuHandler.GetMenu).Methods("GET")
	api.HandleFunc("/menu", menuHandler.AddMenuItem).Methods("POST")
	api.HandleFunc("/menu/{id}", menuHandler.UpdateMenuItem).Methods("PUT")
	api.HandleFunc("/menu/{id}", menuHandler.DeleteMenuItem).Methods("DELETE")

	// Order routes
	api.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	api.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	api.HandleFunc("/orders/{id}/status", orderHandler.UpdateOrderStatus).Methods("PUT")

	// WebSocket endpoint for real-time updates
	r.HandleFunc("/ws", utils.HandleWebSocket)

	// Add a route handler for the root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Food Court API"))
	}).Methods("GET")

	// Wrap router with CORS middleware
	handler := c.Handler(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Using database URL: %s", os.Getenv("DATABASE_URL"))
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Printf("Server failed to start: %v", err)
		log.Fatal(err)
	}
}