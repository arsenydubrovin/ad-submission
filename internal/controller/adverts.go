package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arsenydubrovin/ad-submission/internal/models"
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
		return ctx.JSON(http.StatusBadRequest, envelope{"error": "invalid request payload"})
	}

	advert := &models.Advert{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		PhotoLinks:  input.PhotoLinks,
	}

	// TODO: validate advert

	if err := c.models.Adverts.Insert(advert); err != nil {
		// TODO: custom errResponse()
		return ctx.JSON(http.StatusInternalServerError, envelope{"error": "failed to insert advert"})
	}

	// TODO: custom writeJSON()
	return ctx.JSON(http.StatusCreated, envelope{"advert": advert})
}

func (c *Controller) fetchAdvertHandler(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// TODO: custom errResponse()
		return ctx.JSON(http.StatusNotFound, envelope{"error": "failed to parse id parameter"})
	}

	advert, err := c.models.Adverts.Get(id)
	if err != nil {
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				// TODO: custom errResponse()
				return ctx.JSON(http.StatusNotFound, envelope{"error": "advert not found"})
			default:
				// TODO: custom errResponse()
				return ctx.JSON(http.StatusInternalServerError, envelope{"error": "failed to fetch advert"})
			}
		}
	}

	// TODO: custom writeJSON()
	return ctx.JSON(http.StatusCreated, envelope{"advert": advert})
}

func (c *Controller) listAdvertsHandler(ctx echo.Context) error {
	// TODO: filters (sort and page)

	adverts, err := c.models.Adverts.GetAll()
	if err != nil {
		// TODO: custom errResponse()
		return ctx.JSON(http.StatusInternalServerError, envelope{"error": "failed to list advert"})
	}

	return ctx.JSON(http.StatusOK, envelope{"number": len(adverts), "adverts": adverts})
}
