package controllers

import (
	"net/http"

	"github.com/Rabiann/weather-mailer/services"
	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	WeatherService services.WeatherService
}

func (w WeatherController) GetWeather(ctx *gin.Context) {
	city, ok := ctx.GetQuery("city")
	if !ok {
		ctx.JSON(400, nil)
		return
	}

	weather, err := w.WeatherService.GetWeather(city)
	if err != nil {
		ctx.JSON(400, nil)
	}

	ctx.JSON(http.StatusOK, weather)
}
