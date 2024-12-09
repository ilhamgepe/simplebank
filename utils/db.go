package utils

import "github.com/jackc/pgx/v5/pgtype"

func NewNullableString(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: value, Valid: true}
}
