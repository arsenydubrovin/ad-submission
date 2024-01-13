package controller

import (
	"github.com/arsenydubrovin/ad-submission/src/internal/models"
	echo "github.com/labstack/echo/v4"
)

type Controller struct {
	server *echo.Echo
	models *models.Models
}

func New(e *echo.Echo, models *models.Models) Controller {
	return Controller{
		server: e,
		models: models,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.POST("/advert", c.createAdvertHandler)
	c.server.GET("/advert/:id", c.fetchAdvertHandler)
	c.server.GET("/adverts", c.listAdvertsHandler)
}

func (c *Controller) Serve(port string) error {
	err := c.server.Start(":" + port)
	if err != nil {
		return err
	}
	return nil
}
