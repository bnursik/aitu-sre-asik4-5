package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type order struct {
	ID        int `json:"id,omitempty"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

var db *sql.DB

func main() {
	var err error
	db, err = openDB()
	if err != nil {
		log.Fatalf("database connection failed on startup: %v", err)
	}

	if err := createTable(); err != nil {
		log.Fatalf("table creation failed on startup: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/orders", ordersHandler)
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("order-service listening on :8080")
	if err := http.ListenAndServe(":8080", logRequests(mux)); err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	host := requiredEnv("DB_HOST")
	port := requiredEnv("DB_PORT")
	user := requiredEnv("DB_USER")
	password := requiredEnv("DB_PASSWORD")
	name := requiredEnv("DB_NAME")

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	database, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	for i := 1; i <= 10; i++ {
		if err := database.Ping(); err == nil {
			return database, nil
		}
		log.Printf("waiting for database, attempt %d/10", i)
		time.Sleep(2 * time.Second)
	}

	return database, database.Ping()
}

func createTable() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.Ping() != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "database unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listOrders(w)
	case http.MethodPost:
		createOrder(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func listOrders(w http.ResponseWriter) {
	if db == nil || db.Ping() != nil {
		http.Error(w, "database unavailable", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, user_id, product_id, quantity FROM orders ORDER BY id")
	if err != nil {
		http.Error(w, "failed to list orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := []order{}
	for rows.Next() {
		var o order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity); err != nil {
			http.Error(w, "failed to read order", http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	writeJSON(w, http.StatusOK, orders)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.Ping() != nil {
		http.Error(w, "database unavailable", http.StatusInternalServerError)
		return
	}

	var o order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if o.UserID <= 0 || o.ProductID <= 0 || o.Quantity <= 0 {
		http.Error(w, "user_id, product_id, and quantity must be positive", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO orders (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id",
		o.UserID, o.ProductID, o.Quantity,
	).Scan(&o.ID)
	if err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, o)
}

func requiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("configuration validation failed: missing required environment variable %s", key)
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
