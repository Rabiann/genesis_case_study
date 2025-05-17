package controllers

import (
	"fmt"
	"net/http"

	"github.com/Rabiann/weather-mailer/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionController struct {
	SubscriptionService services.SubscriptionService
	TokenService        services.TokenService
	EmailService        services.MailingService
}

func (s SubscriptionController) Subscribe(ctx *gin.Context) {
	var subscription services.Subscription
	if err := ctx.BindJSON(&subscription); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	id, err := s.SubscriptionService.AddSubscription(s.SubscriptionService.MapSubscription(subscription))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	token, err := s.TokenService.CreateToken(id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	url := fmt.Sprintf("http://localhost:8000/api/confirm/%s", token)

	if err := s.EmailService.SendConfirmationLetter(subscription.Email, url); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.JSON(200, url)
}

func (s SubscriptionController) Confirm(ctx *gin.Context) {
	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please use correct token"})
		return
	}

	if err := s.TokenService.UseToken(token); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is invalid"})
		return
	}

	subscriberId, err := s.TokenService.GetSubscription(token)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := s.SubscriptionService.ActivateSubscription(subscriberId); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (s SubscriptionController) Unsubscribe(ctx *gin.Context) {
	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please use correct token"})
		return
	}

	if err := s.TokenService.UseToken(token); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is invalid"})
		return
	}

	subscriberId, err := s.TokenService.GetSubscription(token)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := s.SubscriptionService.DeleteSubscription(subscriberId); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
