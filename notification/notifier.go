package notification

import (
	"fmt"
	"os"
	"time"

	"github.com/Rabiann/weather-mailer/services"
	"github.com/Rabiann/weather-mailer/services/models"
	"github.com/go-co-op/gocron/v2"
)

const Day = time.Hour * 24

type Period int

const (
	Hourly Period = iota
	Daily
)

type Notifier struct {
	weatherService      services.WeatherService
	subscriptionService services.SubscriptionService
	mailingService      services.MailingService
	tokenService        services.TokenService
}

func NewNotifier(weatherService services.WeatherService, subscriptionService services.SubscriptionService, mailingService services.MailingService, tokenService services.TokenService) Notifier {
	return Notifier{
		weatherService:      weatherService,
		subscriptionService: subscriptionService,
		mailingService:      mailingService,
		tokenService:        tokenService,
	}
}

func (n Notifier) RunNotifier() {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			Day,
		),
		gocron.NewTask(
			n.RunSendingPipeline,
			Daily,
		),
	)

	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			time.Hour,
		),
		gocron.NewTask(
			n.RunSendingPipeline,
			Hourly,
		),
	)

	if err != nil {
		panic(err)
	}

	s.Start()

	// block thread, run scheduler infinitely
	select {}
}

func (n Notifier) RunSendingPipeline(period Period) {
	var subscribers []models.Subscription
	var per string
	var weather services.Weather
	var ok bool
	var err error

	cache := make(map[string]services.Weather)

	if period == Daily {
		per = "daily"
	} else {
		per = "hourly"
	}

	result := n.subscriptionService.Db.Where("frequency = ?", per).Find(&subscribers)
	if result.Error != nil {
		panic(result.Error)
	}

	for _, sub := range subscribers {
		fmt.Println(sub)
		weather, ok = cache[sub.City]
		if !ok {
			weather, err = n.weatherService.GetWeather(sub.City)
			if err != nil {
				panic(err)
			}

			cache[sub.City] = weather
		}
		token, err := n.tokenService.CreateToken(sub.ID)
		if err != nil {
			fmt.Println(err)
		}

		baseUrl := os.Getenv("BASE_URL")
		if os.Getenv("HTTPS") == "1" {
			baseUrl = "https://" + baseUrl
		} else {
			baseUrl = "http://" + baseUrl
		}

		url := fmt.Sprintf("%s/api/unsubscribe/%s", baseUrl, token)

		if err = n.mailingService.SendWeatherReport(sub.Email, per, sub.City, weather, url); err != nil {
			fmt.Println(err)
		}
	}
}
