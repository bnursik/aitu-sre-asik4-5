package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []product{
	{ID: 1, Name: "Mechanical Keyboard", Price: 89.99},
	{ID: 2, Name: "Gaming Mouse", Price: 49.99},
	{ID: 3, Name: "USB-C Hub", Price: 29.99},
	{ID: 4, Name: "Monitor", Price: 199.99},
	{ID: 5, Name: "Headset", Price: 59.99},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/products", productsHandler)
	mux.HandleFunc("/products/", productByIDHandler)
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("product-service listening on :8080")
	if err := http.ListenAndServe(":8080", logRequests(mux)); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func productByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	for _, p := range products {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
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
