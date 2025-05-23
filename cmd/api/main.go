// cmd/api/main.go
package main

import (
	"log"
	"net/http"

	"github.com/Treblle/treblle-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/you/linkedinify/internal/config"
	"github.com/you/linkedinify/internal/router"
)

func main() {
	// Load environment variables
	// Try multiple paths to find the .env file
	err := godotenv.Load(".env", "../../.env", "/Users/pratimbhosale/Desktop/hobbyprojects/linkedinify/.env")
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	} else {
		log.Println("✓ Successfully loaded .env file")
	}
	cfg := config.Load()

	// Initialize Treblle if credentials are provided
	if cfg.TreblleToken != "" && cfg.TreblleAPIKey != "" {
		treblle.Configure(treblle.Configuration{
			SDK_TOKEN: cfg.TreblleToken,
			API_KEY:   cfg.TreblleAPIKey, // Using the API key from environment variable
			Debug:     true,
		})
		log.Println("✓ Treblle monitoring enabled with debug mode")
	} else {
		log.Println("⚠ Treblle monitoring disabled - missing TREBLLE_SDK_TOKEN or TREBLLE_API_KEY")
	}

	// Create the router
	app := router.New(cfg)

	// Create a new Chi router
	r := chi.NewRouter()

	// Apply Treblle middleware using the recommended Chi pattern
	r.Use(treblle.Middleware)

	// Mount the app router under the Chi router
	r.Mount("/", app)

	server := &http.Server{
		Addr:    cfg.HTTPAddr, // Use the default port 8080 from config
		Handler: r,
	}

	log.Printf("⇢ Server starting on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
