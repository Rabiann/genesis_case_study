package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mailersend/mailersend-go"
)

type MailingService struct {
	Client               *mailersend.Mailersend
	ConfirmationTemplate string
	WeatherTemplate      string
}

func NewMailingService() (MailingService, error) {
	var ms MailingService
	api_key := os.Getenv("MAILERSEND_API_KEY")
	if api_key == "" {
		return ms, errors.New("MAILERSEND_API_KEY not set")
	}
	mg := mailersend.NewMailersend(api_key)
	ms.Client = mg

	confirmationTemplate, err := os.ReadFile("./templates/confirmationMail.tmpl")
	if err != nil {
		return ms, err
	}

	weatherTemplate, err := os.ReadFile("./templates/weatherMail.tmpl")
	if err != nil {
		return ms, err
	}

	ms.ConfirmationTemplate = string(confirmationTemplate)
	ms.WeatherTemplate = string(weatherTemplate)
	return ms, nil
}

func (s MailingService) buildConfirmationLetter(email string) string {
	return strings.Replace(s.ConfirmationTemplate, "{}", email, 2)
}

func (s MailingService) buildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string {
	let := strings.Replace(s.WeatherTemplate, "{City}", city, 1)
	let = strings.Replace(let, "{Temperature}", temp, 1)
	let = strings.Replace(let, "{Humidity}", humid, 1)
	let = strings.Replace(let, "{UnsubscribeLink}", unsubscribe, 1)
	let = strings.Replace(let, "{Description}", description, 1)
	return let
}

func (s MailingService) SendConfirmationLetterWithAPI(recipient string, confirmationUrl string) error {
	from := mailersend.From{
		Name:  "Confirmator",
		Email: os.Getenv("SENDER_MAIL"),
	}
	to := []mailersend.Recipient{
		{
			Email: recipient,
		},
	}
	subject := "Confirm Weather Subscription"
	body := s.buildConfirmationLetter(confirmationUrl)

	message := s.Client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(to)
	message.SetSubject(subject)
	message.SetHTML(body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := s.Client.Email.Send(ctx, message)
	return err
}

func (s MailingService) SendWeatherReport(recipient string, period string, city string, weather Weather, unsibscribingUrl string) error {
	from := mailersend.From{
		Name:  "Reporter",
		Email: os.Getenv("SENDER_MAIL"),
	}
	to := []mailersend.Recipient{
		{
			Email: recipient,
		},
	}
	subject := fmt.Sprintf("%s report for %s", period, city)
	body := s.buildWeatherLetter(city, fmt.Sprintf("%.1f", weather.Temperature), fmt.Sprintf("%.1f", weather.Humidity), weather.Description, unsibscribingUrl)

	message := s.Client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(to)
	message.SetSubject(subject)
	message.SetHTML(body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := s.Client.Email.Send(ctx, message)
	return err
}
