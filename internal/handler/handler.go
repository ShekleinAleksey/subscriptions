package handler

import (
	_ "github.com/ShekleinAleksey/subscriptions/docs"
	"github.com/ShekleinAleksey/subscriptions/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	SubscriptionHandler *SubscriptionHandler
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		SubscriptionHandler: NewSubscriptionHandler(s.SubscriptionService),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.GET("", h.SubscriptionHandler.ListSubscriptions)
			subscriptions.POST("", h.SubscriptionHandler.CreateSubscription)
			subscriptions.GET("/:id", h.SubscriptionHandler.GetSubscription)
			subscriptions.PUT("/:id", h.SubscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:id", h.SubscriptionHandler.DeleteSubscription)
			subscriptions.GET("/summary", h.SubscriptionHandler.GetSubscriptionSummary)
		}
	}

	return router
}
