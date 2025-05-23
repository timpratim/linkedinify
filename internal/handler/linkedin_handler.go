// internal/handler/linkedin_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/you/linkedinify/internal/middleware"
	"github.com/you/linkedinify/internal/service"
)

type LinkedInHandler struct {
	svc *service.LinkedInService
}

func NewLinkedIn(svc *service.LinkedInService) *LinkedInHandler { return &LinkedInHandler{svc} }

func (h *LinkedInHandler) Routes(secret []byte) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth(secret))
	r.Post("/", h.transform)
	r.Get("/", h.history)
	return r
}

type reqBody struct {
	Text string `json:"text"`
}

func (h *LinkedInHandler) transform(w http.ResponseWriter, r *http.Request) {
	var in reqBody
	if json.NewDecoder(r.Body).Decode(&in) != nil || in.Text == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uid := middleware.UserID(r.Context())
	out, err := h.svc.Transform(r.Context(), uid, in.Text)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"post": out})
}

func (h *LinkedInHandler) history(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserID(r.Context())
	items, err := h.svc.History(r.Context(), uid)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	type item struct {
		ID    uuid.UUID `json:"id"`
		Input string    `json:"input"`
		Post  string    `json:"post"`
	}
	var res []item
	for _, p := range items {
		res = append(res, item{ID: p.ID, Input: p.InputText, Post: p.OutputText})
	}
	_ = json.NewEncoder(w).Encode(res)
}
