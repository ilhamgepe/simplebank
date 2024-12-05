package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) mount() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// without authentication middleware
	r.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", s.createUser)
			r.Post("/login", s.loginUser)
		})
	})

	// with authentication middleware
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(s.config, s.tokenMaker))
		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", s.createAccount)
			r.Get("/", s.listAccounts)
			r.Get("/{id}", s.getAccount)
			// r.Patch("/{id}", s.updateAccount)
			// r.Delete("/{id}", s.deleteAccount)
		})

		r.Route("/transfers", func(r chi.Router) {
			r.Post("/", s.createTransfer)
			// r.Get("/", s.listTransfers)
		})
	})

	s.router = r
}
