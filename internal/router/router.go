// internal/router/router.go
package router

import (
	"context"
	"log"
	"net/http"

	"github.com/Treblle/treblle-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/you/linkedinify/internal/ai"
	"github.com/you/linkedinify/internal/config"
	"github.com/you/linkedinify/internal/db"
	"github.com/you/linkedinify/internal/handler"
	"github.com/you/linkedinify/internal/repository"
	"github.com/you/linkedinify/internal/service"
)

// treblleSetupMiddleware configures Treblle and returns a middleware that includes tracing.
func treblleSetupMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	// Configure Treblle if credentials are provided
	if cfg.TreblleToken != "" && cfg.TreblleAPIKey != "" {
		treblle.Configure(treblle.Configuration{
			SDK_TOKEN: cfg.TreblleToken,
			API_KEY:   cfg.TreblleAPIKey,
			Debug:     true,
		})
		log.Println("✓ Treblle monitoring enabled")

		// Return a middleware that combines trace and Treblle middleware
		return func(next http.Handler) http.Handler {
			return traceMiddleware(treblle.Middleware(next))
		}
	}

	// If Treblle is not configured, log it and only use the trace middleware
	log.Println("⚠ Treblle monitoring disabled - missing credentials")
	return traceMiddleware
}

func New(cfg config.Config) *chi.Mux {
	database := db.New(cfg)
	userRepo := repository.NewUserRepo(database)
	postRepo := repository.NewPostRepo(database)

	authSvc := service.NewAuth(userRepo, cfg)
	aiClient := ai.NewOpenAI(cfg.OpenAIToken)
	liSvc := service.NewLinkedIn(aiClient, postRepo)

	authH := handler.NewAuth(authSvc)
	liH := handler.NewLinkedIn(liSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5, "gzip"))
	r.Use(treblleSetupMiddleware(cfg))

	// Create API v1 router
	v1Router := chi.NewRouter()
	v1Router.Mount("/auth", authH.Routes())
	v1Router.Mount("/posts", liH.Routes(cfg.JWTSecret))

	// Mount v1 router under /api/v1
	r.Mount("/api/v1", v1Router)

	return r
}

// A private type for context keys to avoid collisions.
type contextKey string

const traceIDKey = contextKey("traceID")

// traceMiddleware generates a unique trace ID and adds it to the request context and response headers.
func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a trace ID is already present in the request header.
		traceID := r.Header.Get("treblle-tag-id")
		if traceID == "" {
			// If not, generate a new one.
			traceID = uuid.New().String()
			// Add it to the request header so subsequent middlewares can see it.
			r.Header.Set("treblle-tag-id", traceID)
		}

		// Also add it to the response header so the client can see which ID was used.
		w.Header().Set("treblle-tag-id", traceID)

		ctx := context.WithValue(r.Context(), traceIDKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
