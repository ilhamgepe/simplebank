package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/token"
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

	authPayload := r.Context().Value(authPayload).(*token.Payload)

	account, err := s.store.GetAccount(r.Context(), int64(idInt))
	if err != nil {
		s.knownSqlError(w, err)
		return
	}

	if account.Owner != authPayload.Username {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  false,
			Data:    nil,
			Message: "Unauthorized",
		})
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
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		tempLimit, err := strconv.Atoi(limit)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Status:  false,
				Message: "Invalid Limit",
				Data:    "Bad Request",
			})
			return
		}
		args.Limit = int32(tempLimit)
	}

	offset := r.URL.Query().Get("offset")
	if offset != "" {
		tempOffset, err := strconv.Atoi(offset)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Status:  false,
				Message: "Invalid Offset",
				Data:    "Bad Request",
			})
			return
		}
		args.Offset = int32(tempOffset)
	}
	authPayload := r.Context().Value(authPayload).(*token.Payload)
	args.Owner = authPayload.Username
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
