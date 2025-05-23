// internal/handler/auth_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/you/linkedinify/internal/service"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuth(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc} }

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.login)
	r.Post("/register", h.register)
	return r
}

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var c creds
	if json.NewDecoder(r.Body).Decode(&c) != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	token, err := h.svc.Login(r.Context(), c.Email, c.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var c creds
	if json.NewDecoder(r.Body).Decode(&c) != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	token, err := h.svc.Register(r.Context(), c.Email, c.Password)
	if err != nil {
		http.Error(w, "conflict", http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
