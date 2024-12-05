package api

import (
	"net/http"
	"time"

	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/utils"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=5,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

type userResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	var err error
	if res, err := s.vStruct(r, &req); err != nil || res != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status: false,
			Data:   res,
		})
		return
	}

	req.Password, err = utils.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Response{
			Status: false,
			Data:   err.Error(),
		})
		return
	}

	user, err := s.store.CreateUser(r.Context(), db.CreateUserParams(req))
	if err != nil {
		s.knownSqlError(w, err)
		return
	}

	rsp := userResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, Response{
		Status:  true,
		Data:    rsp,
		Message: "success",
	})
}

type loginUserRequest struct {
	Username string `json:"username" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        userResponse
}

func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest
	var err error
	if res, err := s.vStruct(r, &req); err != nil || res != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status: false,
			Data:   res,
		})
		return
	}

	user, err := s.store.GetUser(r.Context(), req.Username)
	if err != nil {
		s.knownSqlError(w, err)
		return
	}

	err = utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, Response{
			Status:  false,
			Data:    nil,
			Message: "Unauthorized",
		})
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Response{
			Status: false,
			Data:   err.Error(),
		})
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User: userResponse{
			Username:         user.Username,
			FullName:         user.FullName,
			PasswordChangeAt: user.PasswordChangeAt,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
		},
	}

	writeJSON(w, http.StatusOK, Response{
		Status:  true,
		Data:    rsp,
		Message: "success",
	})
}
