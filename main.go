package main

import (
	"log"
	"net/http"
	"time"
)

// Middleware type as before
type Middleware func(http.Handler) http.Handler

// App struct to hold our routes and middleware
type App struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

// NewApp creates and returns a new App with an initialized ServeMux and middleware slice
func NewApp() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
	}
}

// Use adds middleware to the chain
func (a *App) Use(mw Middleware) {
	a.middlewares = append(a.middlewares, mw)
}

// Handle registers a handler for a specific route, applying all middleware
func (a *App) Handle(pattern string, handler http.Handler) {
	finalHandler := handler
	for i := len(a.middlewares) - 1; i >= 0; i-- {
		finalHandler = a.middlewares[i](finalHandler)
	}
	a.mux.Handle(pattern, finalHandler)
}

// ListenAndServe starts the application server
func (a *App) ListenAndServe(address string) error {
	return http.ListenAndServe(address, a.mux)
}

func main() {
	app := NewApp()

	// Add middleware
	app.Use(LoggingMiddleware)
	app.Use(AnotherMiddleware)

	// Add routes
	app.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	// Start the server
	log.Fatal(app.ListenAndServe(":3000"))
}

// LoggingMiddleware logs the request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

// AnotherMiddleware is an example of additional middleware
func AnotherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do something before
		log.Println("Before handling request")
		next.ServeHTTP(w, r)
		// Do something after
		log.Println("After handling request")
	})
}

// ApplyMiddleware applies multiple middleware to a http.Handler
func ApplyMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
