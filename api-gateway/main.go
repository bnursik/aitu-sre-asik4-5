package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const serviceName = "api-gateway"

type route struct {
	prefix      string
	servicePath string
	target      string
}

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"service", "method", "path", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration)
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

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(data)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
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
		recorder := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		status := recorder.status
		if status == 0 {
			status = http.StatusOK
		}
		log.Printf("%s %s %s", r.Method, r.URL.Path, duration)

		if r.URL.Path == "/metrics" {
			return
		}
		path := routeLabel(r.URL.Path)
		requestsTotal.WithLabelValues(serviceName, r.Method, path, strconv.Itoa(status)).Inc()
		requestDuration.WithLabelValues(serviceName, r.Method, path).Observe(duration.Seconds())
	})
}

func routeLabel(path string) string {
	switch {
	case path == "/health":
		return "/health"
	case path == "/api/auth" || strings.HasPrefix(path, "/api/auth/"):
		return "/api/auth"
	case path == "/api/users" || strings.HasPrefix(path, "/api/users/"):
		return "/api/users"
	case path == "/api/products":
		return "/api/products"
	case strings.HasPrefix(path, "/api/products/"):
		return "/api/products/{id}"
	case path == "/api/orders" || strings.HasPrefix(path, "/api/orders/"):
		return "/api/orders"
	case path == "/api/payments" || strings.HasPrefix(path, "/api/payments/"):
		return "/api/payments"
	case path == "/api/notifications" || strings.HasPrefix(path, "/api/notifications/"):
		return "/api/notifications"
	default:
		return "unknown"
	}
}
