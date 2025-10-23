package config

import (
	"os"
)

type Config struct {
	KafkaBrokers string
	KafkaTopic   string
	MongoURI     string
	MongoDB      string
	RESENDAPIKEY string
}

func Load() Config {
	return Config{
		KafkaBrokers: getenv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:   getenv("KAFKA_TOPIC", "notifications"),
		MongoURI:     getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getenv("MONGO_DB", "notificationdb"),
		RESENDAPIKEY: getenv("RESEND_API_KEY", "re_XGKwg7gZ_2tUeHfr4hBgdT59BAkH5v5J3"),
	}
}

func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
