package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) mount() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", s.createAccount)
		r.Get("/", s.listAccounts)
		r.Get("/{id}", s.getAccount)
		// r.Patch("/{id}", s.updateAccount)
		// r.Delete("/{id}", s.deleteAccount)
	})

	s.router = r
}
