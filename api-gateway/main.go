package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type route struct {
	prefix      string
	servicePath string
	target      string
}

func main() {
	routes := []route{
		{prefix: "/api/auth", servicePath: "", target: "http://auth-service:8080"},
		{prefix: "/api/users", servicePath: "/users", target: "http://user-service:8080"},
		{prefix: "/api/products", servicePath: "/products", target: "http://product-service:8080"},
		{prefix: "/api/orders", servicePath: "/orders", target: "http://order-service:8080"},
		{prefix: "/api/payments", servicePath: "/payments", target: "http://payment-service:8080"},
		{prefix: "/api/notifications", servicePath: "/notifications", target: "http://notification-service:8080"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	for _, r := range routes {
		mux.Handle(r.prefix+"/", proxyHandler(r.prefix, r.servicePath, r.target))
		mux.Handle(r.prefix, proxyHandler(r.prefix, r.servicePath, r.target))
	}

	log.Println("api-gateway listening on :8080")
	if err := http.ListenAndServe(":8080", logRequests(mux)); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func proxyHandler(prefix, servicePath, target string) http.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("invalid target URL %s: %v", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		suffix := strings.TrimPrefix(req.URL.Path, prefix)
		req.URL.Path = servicePath + suffix
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
		req.Host = targetURL.Host
	}

	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
