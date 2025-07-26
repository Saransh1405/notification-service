package main

import (
	"log"
	"notification-service/config"
	"notification-service/handlers"
	"notification-service/kafka"
	"notification-service/mongo"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect MongoDB
	mongo.Connect(cfg)

	// Start Kafka consumer (in goroutine)
	brokers := []string{cfg.KafkaBrokers}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	go kafka.StartConsumers(brokers, saramaConfig)

	// Start Gin HTTP server for health
	r := gin.Default()
	r.GET("/health", handlers.HealthCheck)
	log.Println("Starting HTTP server on :8082...")
	r.Run(":8082")
}
