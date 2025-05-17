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
	Client *mailersend.Mailersend
}

func NewMailingService() (MailingService, error) {
	var ms MailingService
	api_key := os.Getenv("MAILERSEND_API_KEY")
	if api_key == "" {
		return ms, errors.New("MAILERSEND_API_KEY not set")
	}
	mg := mailersend.NewMailersend(api_key)
	ms.Client = mg
	return ms, nil
}

func (s MailingService) SendConfirmationLetter(recipient string, confirmationUrl string) error {
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
	body := fmt.Sprintf("Dear %s, please confirm subscription: %s.", recipient, confirmationUrl)

	message := s.Client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(to)
	message.SetSubject(subject)
	message.SetText(body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := s.Client.Email.Send(ctx, message)
	return err
}

func (s MailingService) SendWeatherReport(recipient string, period string, city string, weather Weather) error {
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
	body := fmt.Sprintf("%s report for %s\nTemperature: %f\nHumidity: %f\nDescription: %s\n\n", strings.Title(period), city, weather.Temperature, weather.Humidity, weather.Description)

	message := s.Client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(to)
	message.SetSubject(subject)
	message.SetText(body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := s.Client.Email.Send(ctx, message)
	return err
}
