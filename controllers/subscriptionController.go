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
	BaseUrl             string
}

func (s SubscriptionController) Subscribe(ctx *gin.Context) {
	var subscription services.Subscription

	subscription.Email = ctx.PostForm("email")
	subscription.City = ctx.PostForm("city")
	subscription.Frequency = ctx.PostForm("period")

	id, err := s.SubscriptionService.AddSubscription(s.SubscriptionService.MapSubscription(subscription))
	if err != nil {
		fmt.Println(err)
		ctx.HTML(409, "alreadysubscribed.html", gin.H{})
		return
	}

	fmt.Println(subscription, id)

	token, err := s.TokenService.CreateToken(id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	url := fmt.Sprintf("%s/api/confirm/%s", s.BaseUrl, token)

	if err := s.EmailService.SendConfirmationLetterWithAPI(subscription.Email, url); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.HTML(http.StatusOK, "needconfirmation.html", gin.H{})
}

func (s SubscriptionController) Confirm(ctx *gin.Context) {
	handleTokenErr := func(ctx *gin.Context, err error, code int) {
		fmt.Println(err)
		ctx.HTML(code, "registrationfailed.html", gin.H{})
	}

	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		handleTokenErr(ctx, err, 400)
		return
	}

	subscriberId, err := s.TokenService.GetSubscription(token)
	if err != nil {
		handleTokenErr(ctx, err, 400)
		return
	}

	if err := s.TokenService.UseToken(token); err != nil {
		handleTokenErr(ctx, err, 404)
		return
	}

	_, err = s.SubscriptionService.ActivateSubscription(subscriberId)
	if err != nil {
		handleTokenErr(ctx, err, 400)
		return
	}

	ctx.HTML(http.StatusOK, "registration.html", gin.H{})
}

func (s SubscriptionController) Unsubscribe(ctx *gin.Context) {
	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please use correct token"})
		return
	}
	subscriberId, err := s.TokenService.GetSubscription(token)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	fmt.Println("subscribet id: ", subscriberId)

	if err := s.TokenService.UseToken(token); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is invalid"})
		return
	}

	if err := s.SubscriptionService.DeleteSubscription(subscriberId); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.HTML(http.StatusOK, "unsubscription.html", gin.H{})
}
