package config

import (
	"os"
	"strconv"
)

type Config struct {
	KafkaBrokers     string
	KafkaTopic       string
	MongoURI         string
	MongoDB          string
	MailtrapHost     string
	MailtrapPort     int
	MailtrapUsername string
	MailtrapPassword string
	FromEmail        string
	FromName         string
}

func Load() Config {
	return Config{
		KafkaBrokers:     getenv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:       getenv("KAFKA_TOPIC", "notifications"),
		MongoURI:         getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:          getenv("MONGO_DB", "notificationdb"),
		MailtrapHost:     getenv("MAILTRAP_HOST", "sandbox.smtp.mailtrap.io"),
		MailtrapPort:     getenvAsInt("MAILTRAP_PORT", 2525),
		MailtrapUsername: getenv("MAILTRAP_USERNAME", "29345f87a2298c"),
		MailtrapPassword: getenv("MAILTRAP_PASSWORD", "1343abbd17b523"),
		FromEmail:        getenv("FROM_EMAIL", "noreply@tribewithvibe.com"),
		FromName:         getenv("FROM_NAME", "Tribe With Vibe"),
	}
}

func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getenvAsInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}
	return fallback
}
