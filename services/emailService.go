package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailingService struct {
	Client               *sendgrid.Client
	ConfirmationTemplate string
	WeatherTemplate      string
}

func NewMailingService() (MailingService, error) {
	var ms MailingService
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	ms.Client = client

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
	return strings.Replace(s.ConfirmationTemplate, "{}", email, 3)
}

func (s MailingService) buildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string {
	let := strings.Replace(s.WeatherTemplate, "{City}", city, 1)
	let = strings.Replace(let, "{Temperature}", temp, 1)
	let = strings.Replace(let, "{Humidity}", humid, 1)
	let = strings.Replace(let, "{UnsubscribeLink}", unsubscribe, 1)
	let = strings.Replace(let, "{Description}", description, 1)
	return let
}

func (s MailingService) SendLetter(from mail.Email, to mail.Email, subject string, content string, ctx context.Context) error {
	message := mail.NewSingleEmail(&from, subject, &to, "", content)
	_, err := s.Client.SendWithContext(ctx, message)
	return err
}

func (s MailingService) SendConfirmationLetter(recipient string, confirmationUrl string) error {
	from := mail.Email{
		Name:    "Confirmator",
		Address: os.Getenv("SENDER_MAIL"),
	}
	to := mail.Email{
		Name:    recipient,
		Address: recipient,
	}

	subject := "Confirm Weather Subscription"
	body := s.buildConfirmationLetter(confirmationUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s.SendLetter(from, to, subject, body, ctx)
	return nil
}

func (s MailingService) SendWeatherReport(recipient string, period string, city string, weather Weather, unsibscribingUrl string) error {
	from := mail.Email{
		Name:    "Reporter",
		Address: os.Getenv("SENDER_MAIL"),
	}
	to := mail.Email{
		Name:    recipient,
		Address: recipient,
	}

	subject := fmt.Sprintf("%s report for %s", period, city)
	body := s.buildWeatherLetter(city, fmt.Sprintf("%.1f", weather.Temperature), fmt.Sprintf("%.1f", weather.Humidity), weather.Description, unsibscribingUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return s.SendLetter(from, to, subject, body, ctx)
}
