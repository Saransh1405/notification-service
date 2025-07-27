package models

type EmailLog struct {
	ID        string                 `json:"id"`
	To        string                 `json:"to"`
	Subject   string                 `json:"subject"`
	Template  string                 `json:"template"`
	Variables map[string]interface{} `json:"variables"`
	Success   bool                   `json:"success"`
	MessageID string                 `json:"messageId"`
	Error     string                 `json:"error"`
	SentAt    int64                  `json:"sentAt"`
	CreatedAt int64                  `json:"createdAt"`
}
