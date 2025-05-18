package services

import (
	"errors"

	"github.com/Rabiann/weather-mailer/services/models"
	"gorm.io/gorm"
)

type Subscription struct {
	Email     string `json:"email"`
	City      string `json:"city"`
	Frequency string `json:"frequency"`
}

type SubscriptionService struct {
	Db *gorm.DB
}

func (s SubscriptionService) MapSubscription(subscriptionRequest Subscription) models.Subscription {
	return models.Subscription{
		Email:     subscriptionRequest.Email,
		Frequency: subscriptionRequest.Frequency,
		City:      subscriptionRequest.City,
		Confirmed: false,
	}
}

func (s SubscriptionService) GetSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := s.Db.Find(&subscriptions)
	return subscriptions, result.Error
}

func (s SubscriptionService) GetSubscriptionById(id uint) (models.Subscription, error) {
	subscription := models.Subscription{ID: id}
	result := s.Db.First(&subscription)
	return subscription, result.Error
}

func (s SubscriptionService) AddSubscription(subscription models.Subscription) (uint, error) {
	if s.Db == nil {
		return 0, nil
	}
	result := s.Db.Create(&subscription)
	return subscription.ID, result.Error
}

func (s SubscriptionService) ActivateSubscription(id uint) (string, error) {
	var subscription models.Subscription
	subscription.ID = id

	result := s.Db.Find(&subscription)
	if result.Error != nil {
		return "", result.Error
	}

	if subscription.Confirmed {
		return "", errors.New("Subscription already confirmed")
	}

	subscription.Confirmed = true
	result = s.Db.Save(subscription)
	return subscription.Email, result.Error
}

func (s SubscriptionService) UpdateSubscription(id uint, new_subscription models.Subscription) error {
	subscription := models.Subscription{ID: id}

	if id != new_subscription.ID {
		return errors.New("IDs differ")
	}

	result := s.Db.Find(&subscription)

	if result.Error != nil {
		return result.Error
	}

	subscription.City = new_subscription.City
	subscription.Confirmed = new_subscription.Confirmed
	subscription.CreatedAt = new_subscription.CreatedAt
	subscription.Email = new_subscription.Email
	subscription.Frequency = new_subscription.Frequency

	result = s.Db.Save(subscription)
	return result.Error
}

func (s SubscriptionService) DeleteSubscription(id uint) error {
	result := s.Db.Delete(&models.Subscription{}, id)
	return result.Error
}

func (s SubscriptionService) Confirm(id uint) error {
	subscription := models.Subscription{ID: id}

	result := s.Db.Find(&subscription)

	if result.Error != nil {
		return result.Error
	}

	subscription.Confirmed = true
	result = s.Db.Save(subscription)
	return result.Error
}
