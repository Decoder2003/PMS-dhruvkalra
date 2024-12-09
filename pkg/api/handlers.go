package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"product-management/pkg/database"
	"product-management/pkg/queue"

	"github.com/gorilla/mux"
	"github.com/lib/pq" // Added for PostgreSQL array support
)

type Product struct {
	ID                 int      `json:"id"`
	UserID             int      `json:"user_id"`
	ProductName        string   `json:"product_name"`
	ProductDescription string   `json:"product_description"`
	ProductImages      []string `json:"product_images"`
	CompressedImages   []string `json:"compressed_product_images"`
	ProductPrice       float64  `json:"product_price"`
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := database.DB.QueryRow(query, product.UserID, product.ProductName, product.ProductDescription,
		pq.Array(product.ProductImages), product.ProductPrice).Scan(&product.ID)

	if err != nil {
		// http.Error(w, "Failed to create product", http.StatusInternalServerError)
		// log.Printf("Database Query Error: %v", err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Send image URLs to RabbitMQ
	queue.PublishToQueue(product.ProductImages)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	userID := r.URL.Query().Get("user_id")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	productName := r.URL.Query().Get("product_name")

	// Base query
	query := `SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price FROM products WHERE 1=1`

	// Dynamically build filters and arguments
	var args []interface{}
	counter := 1

	if userID != "" {
		query += " AND user_id = $" + strconv.Itoa(counter)
		counter++
		args = append(args, userID)
	}
	if minPrice != "" {
		query += " AND product_price >= $" + strconv.Itoa(counter)
		counter++
		args = append(args, minPrice)
	}
	if maxPrice != "" {
		query += " AND product_price <= $" + strconv.Itoa(counter)
		counter++
		args = append(args, maxPrice)
	}
	if productName != "" {
		query += " AND product_name ILIKE $" + strconv.Itoa(counter)
		counter++
		args = append(args, "%"+productName+"%")
	}

	// Execute the query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		log.Printf("Query Error: %v", err)
		return
	}
	defer rows.Close()

	// Parse results
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription, pq.Array(&product.ProductImages), pq.Array(&product.CompressedImages), &product.ProductPrice)
		if err != nil {
			http.Error(w, "Failed to parse products", http.StatusInternalServerError)
			log.Printf("Row Scan Error: %v", err)
			return
		}
		products = append(products, product)
	}

	// Respond with filtered results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product Product
	query := `SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price
              FROM products WHERE id = $1`
	err = database.DB.QueryRow(query, productID).Scan(&product.ID, &product.UserID, &product.ProductName,
		&product.ProductDescription, pq.Array(&product.ProductImages), pq.Array(&product.CompressedImages),
		&product.ProductPrice)

	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}
