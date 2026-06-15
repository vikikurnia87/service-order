package utils

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v5"
)

// Pesan error standar request HTTP — terpusat agar konsisten lintas handler.
const (
	MsgInvalidBody = "invalid request body"
	MsgValidation  = "validation failed"
	MsgInvalidID   = "invalid id"
)

// ParamID mem-parse path param "id" sebagai int64 positif (id <= 0 ditolak).
func ParamID(c *echo.Context) (int64, error) {
	return ParamNamedID(c, "id")
}

// ParamNamedID mem-parse path param bernama `name` sebagai int64 positif.
func ParamNamedID(c *echo.Context, name string) (int64, error) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New(name + " must be a positive integer")
	}
	return id, nil
}
