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
	if err := s.vStruct(w, r, &req); err != nil {
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

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err == nil {
		args.Limit = int32(limit)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err == nil {
		args.Offset = int32(offset)
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