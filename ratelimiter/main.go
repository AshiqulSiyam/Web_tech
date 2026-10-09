package main

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)


type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)


func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if client, exists := clients[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	
	limiter := rate.NewLimiter(1, 3)
	clients[ip] = &Client{limiter: limiter, lastSeen: time.Now()}
	return limiter
}


func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := getVisitor(ip)

		
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests) //429
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ok":      false,
				"message": "Too many requests. Please try again later.",
			})
			return
		}

	
		next.ServeHTTP(w, r)
	})
}


func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome! Request successful.",
	})
}

func main() {
	mux := http.NewServeMux()

	
	mux.Handle("/api/v1/data", rateLimitMiddleware(http.HandlerFunc(apiHandler)))

	addr := ":8080"
	println("Server running on", addr)
	http.ListenAndServe(addr, mux)
}








// 1..10 | ForEach-Object {
// >>     Write-Host "Request ${_}:"
// >>     curl.exe -s http://localhost:8080/api/v1/data
// >>     Write-Host ""
// >> }