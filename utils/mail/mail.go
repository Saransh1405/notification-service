package mail

import (
	"encoding/json"
	"fmt"
	"log"
	"notification-service/config"
	"notification-service/models"
	"notification-service/mongo"
	"time"
)

var emailService *EmailService

func IntializeEmailService(config config.Config) {
	emailService = NewEmailService(config)
	log.Println("Email service initialized")
}

func HandleMail(notification map[string]interface{}) error {
	log.Printf("Processing mail notification: %+v", notification)

	// Extract email data from notification
	emailData, err := extractEmailData(notification)
	if err != nil {
		log.Printf("Error extracting email data: %v", err)
		return fmt.Errorf("failed to extract email data: %w", err)
	}

	log.Printf("Extracted email data: To=%s, Subject=%s, Template=%s",
		emailData.To, emailData.Subject, emailData.Template)

	// Send email
	result, err := emailService.SendEmail(emailData)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log email result to MongoDB
	emailLog := models.EmailLog{
		To:        emailData.To,
		Subject:   emailData.Subject,
		Template:  emailData.Template,
		Variables: emailData.Variables,
		Success:   result.Success,
		MessageID: result.MessageID,
		Error:     result.Error,
		SentAt:    time.Now().UnixMilli(),
		CreatedAt: time.Now().UnixMilli(),
	}

	if err := mongo.SaveEmailLog(emailLog); err != nil {
		log.Printf("Error saving email log: %v", err)
		// Don't return error here as email was sent successfully
	}

	if result.Success {
		log.Printf("Email sent successfully to %s (MessageID: %s)", emailData.To, result.MessageID)
	} else {
		log.Printf("Email failed to send to %s: %s", emailData.To, result.Error)
	}

	return nil
}

// extractEmailData extracts and validates email data from notification
func extractEmailData(notification map[string]interface{}) (models.EmailData, error) {
	// Convert notification data to JSON and back to ensure proper structure
	dataBytes, err := json.Marshal(notification)
	if err != nil {
		return models.EmailData{}, fmt.Errorf("failed to marshal notification data: %w", err)
	}

	var emailData models.EmailData
	if err := json.Unmarshal(dataBytes, &emailData); err != nil {
		return models.EmailData{}, fmt.Errorf("failed to unmarshal email data: %w", err)
	}

	// Additional validation
	if emailData.To == "" {
		return models.EmailData{}, fmt.Errorf("recipient email is required")
	}
	if emailData.Subject == "" {
		return models.EmailData{}, fmt.Errorf("email subject is required")
	}
	if emailData.Template == "" {
		return models.EmailData{}, fmt.Errorf("email template is required")
	}

	return emailData, nil
}
