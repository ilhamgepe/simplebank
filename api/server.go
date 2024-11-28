package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
)

type Server struct {
	store    *db.Store
	router   *chi.Mux
	validate *validator.Validate
}

func NewServer(store *db.Store, validator *validator.Validate) *Server {
	server := &Server{
		store:    store,
		validate: validator,
	}
	server.mount()

	return server
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

type Response struct {
	Status  bool   `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
	}
}
