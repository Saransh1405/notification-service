package mongo

import (
	"context"
	"log"
	"notification-service/config"
	"notification-service/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Connect(cfg config.Config) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Mongo new client error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Mongo connect error: %v", err)
	}
	db = client.Database(cfg.MongoDB)
	log.Printf("MongoDB connected: %s", cfg.MongoDB)
}

func SaveNotification(notif models.Notification) error {
	if db == nil {
		return nil // skip if not connected
	}
	_, err := db.Collection("notifications").InsertOne(context.TODO(), notif)
	return err
}
