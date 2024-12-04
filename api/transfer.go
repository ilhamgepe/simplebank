package api

import (
	"fmt"
	"net/http"

	db "github.com/ilhamgepe/simplebank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Currency      string `json:"currency" validate:"required,oneof=USD EUR CAD"`
}

func (s *Server) createTransfer(w http.ResponseWriter, r *http.Request) {
	var req transferRequest
	if res, err := s.vStruct(r, &req); err != nil || res != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  false,
			Data:    res,
			Message: "Bad Request",
		})
		return
	}

	if !s.validAccount(w, r, req.FromAccountID, req.Currency) {
		return
	}

	account, err := s.store.TransferTx(r.Context(), db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
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

func (s *Server) validAccount(w http.ResponseWriter, r *http.Request, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(r.Context(), accountID)
	if err != nil {
		s.knownSqlError(w, err)
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  false,
			Message: err.Error(),
		})
		return false
	}
	return true
}
