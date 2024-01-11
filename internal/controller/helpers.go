package controller

import (
	"log/slog"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type envelope map[string]any

func (c *Controller) errorResponse(ctx echo.Context, status int, mgs any) error {
	return ctx.JSON(status, envelope{"error": mgs})
}

func (c *Controller) serverErrorResponse(ctx echo.Context, err error) error {
	slog.Error("internal server error", slog.Any("error", err))

	msg := "the server has encountered an error and cannot process the request"
	return c.errorResponse(ctx, http.StatusInternalServerError, msg)
}

func (c *Controller) badRequestResponse(ctx echo.Context, err error) error {
	return c.errorResponse(ctx, http.StatusBadRequest, err)
}

func (c *Controller) notFoundResponse(ctx echo.Context) error {
	msg := "the requested resource was not found"
	return c.errorResponse(ctx, http.StatusNotFound, msg)
}

func (c *Controller) failedValidationResponse(ctx echo.Context, errors map[string]string) error {
	return c.errorResponse(ctx, http.StatusUnprocessableEntity, errors)
}
