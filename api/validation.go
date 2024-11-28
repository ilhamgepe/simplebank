package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Server) knownSqlError(w http.ResponseWriter, err error) {
	// Cek apakah error merupakan pgconn.PgError
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // Duplicate key violation
			// Duplikasi entri (contoh: duplikasi unique key)
			writeJSON(w, http.StatusConflict, Response{
				Status:  false,
				Message: "Duplicate entry detected",
			})
		case "23503": // Foreign key violation
			// Pelanggaran foreign key
			writeJSON(w, http.StatusBadRequest, Response{
				Status:  false,
				Message: "Foreign key constraint violation",
			})
		case "23514": // Check constraint violation
			// Pelanggaran check constraint
			writeJSON(w, http.StatusBadRequest, Response{
				Status:  false,
				Message: "Check constraint violation",
			})
		default:
			// Menangani error lainnya
			writeJSON(w, http.StatusInternalServerError, Response{
				Status:  false,
				Message: "Database error occurred: " + pgErr.Message,
			})
		}
	} else {
		if errors.Is(err, pgx.ErrNoRows) {
			// Tangani kasus tidak ada baris yang ditemukan
			writeJSON(w, http.StatusNotFound, Response{
				Status:  false,
				Message: "No data found",
			})
			return
		}
		// Jika bukan pgconn.PgError
		log.Println("Non-pgx error:", err)
		writeJSON(w, http.StatusInternalServerError, Response{
			Status:  false,
			Message: "Internal server error",
		})

	}
}

func (s *Server) vStruct(w http.ResponseWriter, r *http.Request, data any) error {
	// baca json dari request
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  false,
			Message: "Validation failed",
			Data:    "Bad Request",
		})
		return err
	}

	// Validasi struct
	err := s.validate.Struct(data)
	if err != nil {
		// Ambil error validasi dan format pesan yang mudah dimengerti
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(e))
		}

		// Kirim response dengan error validasi
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  false,
			Message: "Validation failed",
			Data:    validationErrors,
		})
		return err
	}
	return nil
}

// Format error validasi agar lebih mudah dibaca
func formatValidationError(e validator.FieldError) string {
	field := e.Field() // Menggunakan Field() untuk mendapatkan nama field
	tag := e.Tag()

	// Tentukan pesan berdasarkan tag validasi
	var message string
	switch tag {
	case "required":
		message = field + " field is required"
	case "gte":
		message = field + " field must be greater than or equal to " + e.Param()
	case "lte":
		message = field + " field must be less than or equal to " + e.Param()
	case "email":
		message = field + " must be a valid email"
	case "oneof":
		message = field + " must be one of " + e.Param()
	default:
		message = field + " is invalid (" + tag + ")"
	}

	return message
}
