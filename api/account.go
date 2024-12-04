package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=USD EUR CAD"`
}

func (s *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if res, err := s.vStruct(r, &req); err != nil || res != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  false,
			Data:    res,
			Message: "Bad Request",
		})
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(r.Context(), arg)
	if err != nil {
		s.knownSqlError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Status:  true,
		Data:    account,
		Message: "success",
	})
}

func (s *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 1 {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  false,
			Message: "Invalid ID",
			Data:    "Bad Request",
		})
		return
	}

	account, err := s.store.GetAccount(r.Context(), int64(idInt))
	if err != nil {
		s.knownSqlError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Status:  true,
		Data:    account,
		Message: "success",
	})
}

func (s *Server) listAccounts(w http.ResponseWriter, r *http.Request) {
	args := db.GetAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := s.store.GetAccounts(r.Context(), args)
	if err != nil {
		s.knownSqlError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, Response{
		Status:  true,
		Data:    accounts,
		Message: "success",
	})
}
