package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"notification-service/config"
	"notification-service/models"
	"strings"
	"text/template"
	"time"
)

type EmailService struct {
	Config   config.Config
	Template map[string]*template.Template
}

func NewEmailService(config config.Config) *EmailService {
	service := &EmailService{
		Config:   config,
		Template: make(map[string]*template.Template),
	}

	// intialize templates
	service.initializeTemplates()

	return service
}

// initializeTemplates initializes the templates
func (s *EmailService) initializeTemplates() {
	// welcome email template
	welcomeTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Welcome</title>
	</head>
	<body>
		<h1>Welcome {{.Name}}!</h1>
		<p>Thank you for joining our service.</p>
		<p>Click <a href="{{.ActivationLink}}">here</a> to activate your account.</p>
		<p>Best regards,<br>{{.CompanyName}}</p>
	</body>
	</html>`

	// Password reset template
	resetTemplate := `
	 <!DOCTYPE html>
	 <html>
	 <head>
		 <meta charset="UTF-8">
		 <title>Password Reset</title>
	 </head>
	 <body>
		 <h1>Password Reset Request</h1>
		 <p>Hello {{.Name}},</p>
		 <p>You requested a password reset. Click <a href="{{.ResetLink}}">here</a> to reset your password.</p>
		 <p>If you didn't request this, please ignore this email.</p>
		 <p>Best regards,<br>{{.CompanyName}}</p>
	 </body>
	 </html>`

	// Notification template
	notificationTemplate := `
	  <!DOCTYPE html>
	  <html>
	  <head>
		  <meta charset="UTF-8">
		  <title>{{.Title}}</title>
	  </head>
	  <body>
		  <h1>{{.Title}}</h1>
		  <p>{{.Message}}</p>
		  <p>Best regards,<br>{{.CompanyName}}</p>
	  </body>
	  </html>`

	// Parse templates
	s.Template["welcome"] = template.Must(template.New("welcome").Parse(welcomeTemplate))
	s.Template["password_reset"] = template.Must(template.New("password_reset").Parse(resetTemplate))
	s.Template["notification"] = template.Must(template.New("notification").Parse(notificationTemplate))
}

// SendEmail sends an email using the provided data
func (s *EmailService) SendEmail(emailData models.EmailData) (*models.EmailResult, error) {
	// Validate email data
	if err := s.validateEmailData(emailData); err != nil {
		return &models.EmailResult{
			Success: false,
			Error:   err.Error(),
			SentAt:  time.Now().UnixMilli(),
		}, err
	}

	// create an SMTP authentication
	auth := smtp.PlainAuth("", s.Config.MailtrapUsername, s.Config.MailtrapPassword, s.Config.MailtrapHost)

	// Build email message
	message, err := s.buildEmailMessage(emailData)
	if err != nil {
		return &models.EmailResult{
			Success: false,
			Error:   err.Error(),
			SentAt:  time.Now().UnixMilli(),
		}, err
	}

	// Send email
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Config.MailtrapHost, s.Config.MailtrapPort),
		auth,
		s.Config.FromEmail,
		[]string{emailData.To},
		[]byte(message),
	)

	if err != nil {
		return &models.EmailResult{
			Success: false,
			Error:   err.Error(),
			SentAt:  time.Now().UnixMilli(),
		}, err
	}

	return &models.EmailResult{
		Success:   true,
		MessageID: s.generateMessageID(),
		SentAt:    time.Now().UnixMilli(),
	}, nil
}

// validateEmailData validates the email data before sending
func (s *EmailService) validateEmailData(emailData models.EmailData) error {
	if emailData.To == "" {
		return fmt.Errorf("recipient email is required")
	}
	if emailData.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if emailData.Template == "" {
		return fmt.Errorf("email template is required")
	}

	// Check if template exists
	if _, exists := s.Template[emailData.Template]; !exists {
		return fmt.Errorf("template '%s' not found", emailData.Template)
	}

	return nil
}

// buildEmailMessage builds the complete email message
func (s *EmailService) buildEmailMessage(emailData models.EmailData) (string, error) {
	// Build email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.Config.FromName, s.Config.FromEmail)
	headers["To"] = emailData.To
	headers["Subject"] = emailData.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	headers["Message-ID"] = fmt.Sprintf("<%s@%s>", s.generateMessageID(), s.Config.FromEmail)

	// Build headers string
	var headerLines []string
	for key, value := range headers {
		headerLines = append(headerLines, fmt.Sprintf("%s: %s", key, value))
	}

	// Get email body from template
	body, err := s.renderTemplate(emailData.Template, emailData.Variables)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return strings.Join(headerLines, "\r\n") + "\r\n\r\n" + body, nil
}

// renderTemplate renders the email template with variables
func (s *EmailService) renderTemplate(templateName string, variables map[string]interface{}) (string, error) {
	tmpl, exists := s.Template[templateName]
	if !exists {
		return "", fmt.Errorf("template '%s' not found", templateName)
	}

	// Add company name to variables
	if variables == nil {
		variables = make(map[string]interface{})
	}
	variables["CompanyName"] = s.Config.FromName

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateMessageID generates a unique message ID
func (s *EmailService) generateMessageID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
