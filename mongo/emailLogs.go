package mongo

import (
	"context"
	"notification-service/models"
)

const MongoNotificationLogsCollection = "notification_logs"

func SaveEmailLog(emailLog models.EmailLog) error {
	collection := db.Collection(MongoNotificationLogsCollection)
	_, err := collection.InsertOne(context.Background(), emailLog)
	return err
}
