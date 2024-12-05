package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/token"
	"github.com/ilhamgepe/simplebank/utils"
)

type Server struct {
	store      db.Store
	router     *chi.Mux
	validate   *validator.Validate
	tokenMaker token.Maker
	config     utils.Config
}

func NewServer(store db.Store, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %v", err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("currency", validCurrency)
	server := &Server{
		store:      store,
		validate:   validate,
		tokenMaker: tokenMaker,
		config:     config,
	}
	server.mount()

	return server, nil
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
