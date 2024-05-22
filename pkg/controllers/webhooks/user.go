package webhooks

import (
	email "client_task/pkg/utils/email"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	db "client_task/pkg/common/db/sqlc"

	"github.com/gin-gonic/gin"
)

type WebhookController struct {
	db        *db.Queries
	ctx       context.Context
	emailLock sync.Mutex // Mutex for email sending
}

func NewWebhookController(db *db.Queries, ctx context.Context) *WebhookController {
	return &WebhookController{db, ctx, sync.Mutex{}}
}

func sendMail(to, subject, body string) error {
	config, err := email.LoadEmailConfig()
	if err != nil {
		return fmt.Errorf("failed to load email config: %w", err)
	}

	err = email.SendEmail(config, to, subject, body)
	if err != nil {
		return fmt.Errorf("error sending email notification: %w", err)
	}
	return nil
}

type WebhookNotification struct {
	EventType string                 `json:"eventType"`
	JobID     string                 `json:"jobId"`
	Data      map[string]interface{} `json:"data"`
}

func (cc *WebhookController) HandleWebhook(ctx *gin.Context) {
	var notification WebhookNotification
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	switch notification.EventType {
	case "job_application":
		if err := cc.notifyJobApplication(notification.Data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "daily_jobs":
		if err := sendDailyJobNotifications(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown event type"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
}

func (cc *WebhookController) notifyJobApplication(data map[string]interface{}) error {
	jobId, ok := data["jobId"].(string)
	if !ok {
		return errors.New("invalid or missing jobId")
	}

	applicant, ok := data["applicant"].(string)
	if !ok {
		return errors.New("invalid or missing applicant")
	}

	skills, ok := data["skills"].([]interface{})
	if !ok {
		return errors.New("invalid or missing skills")
	}

	for _, skill := range skills {
		skillName, ok := skill.(string)
		if !ok {
			return errors.New("invalid skill format")
		}

		matchingUsers, err := cc.db.GetAllUserBySkillsName(cc.ctx, skillName)
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}

		for _, user := range matchingUsers {
			go cc.sendNotificationEmail(user.Email, jobId, applicant)
		}
	}

	return nil
}

func (cc *WebhookController) sendNotificationEmail(to, jobId, applicant string) {
	cc.emailLock.Lock()
	defer cc.emailLock.Unlock()

	subject := "New Job Application for " + jobId
	body := fmt.Sprintf("A new job application has been received from %s for the job %s.", applicant, jobId)

	if err := sendMail(to, subject, body); err != nil {
		log.Printf("Error sending email notification: %v", err)
	}
}

func sendDailyJobNotifications() error {
	return nil
}
