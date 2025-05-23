// internal/router/router.go
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/you/linkedinify/internal/ai"
	"github.com/you/linkedinify/internal/config"
	"github.com/you/linkedinify/internal/db"
	"github.com/you/linkedinify/internal/handler"
	"github.com/you/linkedinify/internal/repository"
	"github.com/you/linkedinify/internal/service"
)

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
	r.Mount("/auth", authH.Routes())
	r.Mount("/linkedinify", liH.Routes(cfg.JWTSecret))
	return r
}
