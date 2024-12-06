package api

import (
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req renewAccessTokenRequest
	var err error
	if res, err := s.vStruct(r, &req); err != nil || res != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status: false,
			Data:   res,
		})
		return
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  false,
			Data:    nil,
			Message: "Unauthorized",
		})
		return
	}

	session, err := s.store.GetSession(r.Context(), refreshPayload.ID)
	if err != nil {
		s.knownSqlError(w, err)
		return
	}

	if session.IsBlocked {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  false,
			Data:    nil,
			Message: "Unauthorized",
		})
		return
	}

	if session.Username != refreshPayload.Username {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  false,
			Data:    nil,
			Message: "Unauthorized",
		})
		return
	}

	if session.RefreshToken != req.RefreshToken {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  false,
			Data:    nil,
			Message: "Unauthorized",
		})
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  false,
			Data:    nil,
			Message: "Expired session",
		})
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(session.Username, s.config.AccessTokenDuration)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Response{
			Status: false,
			Data:   err.Error(),
		})
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}

	writeJSON(w, http.StatusOK, Response{
		Status:  true,
		Data:    rsp,
		Message: "success",
	})
}
