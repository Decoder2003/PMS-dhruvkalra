package main

import (
	"log"
	"net/http"

	"product-management/pkg/api"
	"product-management/pkg/cache"
	"product-management/pkg/database"
	"product-management/pkg/logger"
	"product-management/pkg/queue"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize components
	database.InitDB()
	queue.InitQueue()
	cache.InitCache()
	logger.InitLogger()

	// API Routes
	r := mux.NewRouter()
	r.HandleFunc("/products", api.CreateProductHandler).Methods("POST")
	r.HandleFunc("/products", api.GetAllProductsHandler).Methods("GET")
	r.HandleFunc("/products/{id}", api.GetProductHandler).Methods("GET")

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
