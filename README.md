# **Product Management System**

This is a backend application for a **Product Management System** with asynchronous image processing. The application is built using **Go** and features RabbitMQ for message queuing, PostgreSQL for data storage, and local storage for compressed images.

## **Features**

### **Product Management API**:

- Add, retrieve, and filter products.
- Supports price range and product name filtering.

### **Asynchronous Image Processing**:

- Upload product images and process them asynchronously.
- Compress and store images locally.
- Update the database with paths to compressed images.

### **Scalable Architecture**:

- Uses RabbitMQ for message queuing.
- Modular design with separate packages for database and queue handling.

---

## **Project Structure**

```plaintext
product-management/
├── cmd/
│   └── main.go            # Entry point for the application
├── pkg/
│   ├── api/
│   │   ├── handlers.go    # HTTP API handlers
│   ├── database/
│   │   ├── db.go          # PostgreSQL connection setup
│   ├── queue/
│       ├── queue.go       # RabbitMQ connection and publishing logic
├── README.md              # Project documentation
├── go.mod                 # Go module dependencies
└── go.sum                 # Checksums for module dependencies
```

## **Installation**

### **Prerequisites**

- Go 1.18 or later
- PostgreSQL installed and running
- RabbitMQ installed and running (or use Docker)
- curl or Postman for API testing

## **Setup Instructions**

**1. Clone the Repository**

```bash
git clone https://github.com/your-repo/product-management.git
cd product-management
```

**2. Install Dependencies**

```bash
go mod tidy
```

**3. Set Up PostgreSQL**

```sql
Create a new database:
postgres
```

```sql
CREATE DATABASE productdb;
Create tables:
sql
Copy code
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_images TEXT[], -- Array of image URLs
    compressed_product_images TEXT[], -- Array of compressed URLs
    product_price DECIMAL(10, 2)
);
```

**4. Start RabbitMQ**

Start RabbitMQ using Docker or your local installation:

```bash copycode
docker run -d --name rabbitmq -p 5672:5672 rabbitmq:management
```

## **Running the Application**

**Start the Application**

```bash
go run cmd/main.go
```

**API Endpoints**

```table
| Method | Endpoint         | Description                                |
|--------|------------------|--------------------------------------------|
| POST   | `/products`      | Create a new product                       |
| GET    | `/products`      | Retrieve all products with optional filters|
| GET    | `/products/{id}` | Retrieve a product by ID                   |
| POST   | `/process-images`| Process images asynchronously              |
```

## **Testing the Endpoints**

**1. Create a Product**

```bash
curl -X POST http://localhost:8080/products \
-H "Content-Type: application/json" \
-d '{
    "user_id": 1,
    "product_name": "Test Product",
    "product_description": "A sample product.",
    "product_images": ["http://example.com/image1.jpg", "http://example.com/image2.jpg"],
    "product_price": 49.99
}'
```

**2. Retrieve All Products**

```bash
curl "http://localhost:8080/products?user_id=1&min_price=10&max_price=50&product_name=Sample"
```

**3. Retrieve a Product by ID**

```bash
curl "http://localhost:8080/products/1"
```

## **Asynchronous Image Processing**

### **How It Works:**

- Images are added to a RabbitMQ queue (image_processing) upon product creation.
- A microservice listens for messages, downloads images, compresses them, and stores the compressed versions locally.
- The compressed_product_images field in the database is updated with the paths of the compressed images.

## **Project Details**

### **Technologies Used**

- Language: Go
- Database: PostgreSQL
- Message Queue: RabbitMQ
- Local Storage: File system for image compression

## **Future Enhancements**

- Add authentication and authorization for the API.
- Implement retry mechanisms for failed image processing tasks.
- Use cloud storage (e.g., S3) for compressed images.
- Implement a dead-letter queue for unprocessed messages.

## **Developer**

- Full Name - Dhruv Kalra
- College - SRM Institute of Science & Technology
- Reg No - RA2111003030194
- Course - B.Tech
- Branch - Computer Science & Engineering
