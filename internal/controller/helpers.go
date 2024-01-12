package controller

import (
	"strconv"
	"strings"

	"github.com/arsenydubrovin/ad-submission/internal/validator"
	echo "github.com/labstack/echo/v4"
)

type envelope map[string]any

func readParamInt(ctx echo.Context, key string, defaultValue int, v *validator.Validator) int {
	p := ctx.QueryParam(key)

	if p == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(p)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}

	return i
}

func readParamString(ctx echo.Context, key string, defaultValue string) string {
	p := ctx.QueryParam(key)

	if p == "" {
		return defaultValue
	}

	return p
}

func readParamCSV(ctx echo.Context, key string, defaultValue []string) []string {
	p := ctx.QueryParam(key)

	if p == "" {
		return defaultValue
	}

	return strings.Split(p, ",")
}
