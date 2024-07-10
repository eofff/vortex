package controllers

import (
	"log"
	"net/http"
	"vortex/services"

	"github.com/labstack/echo/v4"
)

type ClientController struct {
	clientService services.IClientService
}

func (c *ClientController) Init(service services.IClientService, e *echo.Echo) {
	c.clientService = service

	e.POST("/", c.Add)
	e.PUT("/", c.Update)
	e.DELETE("/", c.Delete)
	e.GET("/", c.UpdateStatuses)
}

func (c *ClientController) Add(ctx echo.Context) error {
	var client services.Client
	err := ctx.Bind(&client)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	err = c.clientService.Add(client)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	return ctx.JSON(http.StatusOK, client)
}

func (c *ClientController) Update(ctx echo.Context) error {
	var client services.Client
	err := ctx.Bind(&client)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	err = c.clientService.Update(client)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	return ctx.JSON(http.StatusOK, client)
}

func (c *ClientController) Delete(ctx echo.Context) error {
	var client services.Client
	err := ctx.Bind(&client)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	err = c.clientService.Delete(client)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	return ctx.JSON(http.StatusOK, client)
}

func (c *ClientController) UpdateStatuses(ctx echo.Context) error {
	err := c.clientService.UpdateAlgorithmStatus()
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	return ctx.JSON(http.StatusOK, nil)
}
