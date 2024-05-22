package router

import (
	"client_task/pkg/controllers/webhooks"

	"github.com/gin-gonic/gin"
)

type WebhookRoutes struct {
	webhookController *webhooks.WebhookController
}

func NewWebhookRoutes(webhookController *webhooks.WebhookController) *WebhookRoutes {
	return &WebhookRoutes{webhookController}
}

func (wr *WebhookRoutes) RegisterRoutes(rg *gin.RouterGroup) {
	webhookGroup := rg.Group("webhook")
	webhookGroup.POST("/", wr.webhookController.HandleWebhook)
}
