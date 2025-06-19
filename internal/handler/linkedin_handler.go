// internal/handler/linkedin_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/you/linkedinify/internal/middleware"
	"github.com/you/linkedinify/internal/service"
)

type LinkedInHandler struct {
	svc service.LinkedInServiceInteractor
}

func NewLinkedIn(svc service.LinkedInServiceInteractor) *LinkedInHandler {
	return &LinkedInHandler{svc: svc}
}

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
		http.Error(w, "bad request: missing text", http.StatusBadRequest)
		return
	}

	p := bluemonday.StrictPolicy()
	sanitizedText := p.Sanitize(in.Text)
	uid := middleware.UserID(r.Context())
	out, err := h.svc.Transform(r.Context(), uid, sanitizedText)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"post": out})
}

func (h *LinkedInHandler) history(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserID(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	items, err := h.svc.History(r.Context(), uid, page, pageSize)
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
