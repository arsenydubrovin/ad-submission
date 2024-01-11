package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arsenydubrovin/ad-submission/internal/models"
	"github.com/arsenydubrovin/ad-submission/internal/validator"
	echo "github.com/labstack/echo/v4"
)

func (c *Controller) createAdvertHandler(ctx echo.Context) error {
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       int      `json:"price"`
		PhotoLinks  []string `json:"photoLinks"`
	}

	if err := ctx.Bind(&input); err != nil {
		return c.badRequestResponse(ctx, err)
	}

	advert := &models.Advert{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		PhotoLinks:  input.PhotoLinks,
	}

	v := validator.New()

	if models.ValidateAdvert(v, advert); !v.Valid() {
		return c.failedValidationResponse(ctx, v.Errors)
	}

	if err := c.models.Adverts.Insert(advert); err != nil {
		return c.serverErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, envelope{"advert": advert})
}

func (c *Controller) fetchAdvertHandler(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return c.notFoundResponse(ctx)
	}

	advert, err := c.models.Adverts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			return c.notFoundResponse(ctx)
		default:
			return c.serverErrorResponse(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, envelope{"advert": advert})
}

func (c *Controller) listAdvertsHandler(ctx echo.Context) error {
	// TODO:
	// add filters (sort and page)
	// add filters validation

	adverts, err := c.models.Adverts.GetAll()
	if err != nil {
		return c.serverErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, envelope{"number": len(adverts), "adverts": adverts})
}
