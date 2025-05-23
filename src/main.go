package main

import (
	"net/http"
	"os"

	"github.com/Rabiann/weather-mailer/controllers"
	"github.com/Rabiann/weather-mailer/notification"
	"github.com/Rabiann/weather-mailer/services"
	"github.com/Rabiann/weather-mailer/services/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		panic("WEATHER_API_KEY is not set")
	}
	addr := os.Getenv("WEATHER_API_ADDR")
	if key == "" {
		panic("WEATHER_API_ADDR is not set")
	}
	base_url := os.Getenv("BASE_URL")
	hdr := os.Getenv("HTTPS")
	if hdr == "1" {
		base_url = "https://" + base_url
	} else {
		base_url = "http://" + base_url
	}

	db := models.ConnectToDatabase()

	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Token{}); err != nil {
		panic(err)
	}

	weatherService := services.WeatherService{Address: addr, Key: key}
	subscriptionService := services.SubscriptionService{Db: db}
	tokenService := services.TokenService{Db: db}
	emailService, err := services.NewMailingService()
	if err != nil {
		panic(err)
	}

	notifier := notification.NewNotifier(weatherService, subscriptionService, emailService, tokenService)
	go notifier.RunNotifier()

	weatherController := controllers.WeatherController{WeatherService: weatherService}
	subscriptionController := controllers.SubscriptionController{SubscriptionService: subscriptionService, TokenService: tokenService, EmailService: emailService, BaseUrl: base_url}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/favicon.ico", "./static/weather.ico")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "subscriptions.html", gin.H{})
	})

	api := router.Group("/api")
	{
		api.GET("/weather", weatherController.GetWeather)
		api.POST("/subscribe", subscriptionController.Subscribe)
		api.GET("/confirm/:token", subscriptionController.Confirm)
		api.GET("/unsubscribe/:token", subscriptionController.Unsubscribe)
	}

	router.Run(":8000")
}
