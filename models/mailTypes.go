package models

// EmailData represents the structure of email data from Kafka
type EmailData struct {
	To          string                 `json:"to"`
	Subject     string                 `json:"subject"`
	Template    string                 `json:"template"`
	Variables   map[string]interface{} `json:"variables"`
	Priority    string                 `json:"priority,omitempty"`
	Attachments []Attachment           `json:"attachments,omitempty"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

// EmailResult represents the result of email sending
type EmailResult struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
	SentAt    int64  `json:"sentAt"`
}
