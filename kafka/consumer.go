package kafka

import (
	"encoding/json"
	"fmt"
	"notification-service/models"
	"notification-service/utils/mail"
	"notification-service/utils/push"
	"notification-service/utils/sms"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	NotificationsTopic = "notifications"
)

func StartConsumers(brokers []string, config *sarama.Config) error {
	go ConsumerNotificationForMail(brokers, config)
	go ConsumerNotificationForSMS(brokers, config)
	go ConsumerNotificationForPush(brokers, config)
	return nil
}

func ConsumerNotificationForMail(brokers []string, config *sarama.Config) error {
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("failed to create consumer", zap.Error(err))
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(NotificationsTopic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("failed to create partition consumer", zap.Error(err))
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	const maxRetries = 3
	var event models.Notification

	for {
		select {
		case err := <-partitionConsumer.Errors():
			fmt.Println("error from partition consumer", zap.Error(err))
		case msg := <-partitionConsumer.Messages():
			var lastError error
			success := false

			for attempt := 0; attempt < maxRetries; attempt++ {
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				notificationData, err := json.Marshal(event)
				if err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				var notification models.Notification
				if err := json.Unmarshal(notificationData, &notification); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				if notification.Type == "mail" {
					notificationData, err := json.Marshal(notification.Data)
					if err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}

					var data map[string]interface{}
					if err := json.Unmarshal(notificationData, &data); err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}

					if err := mail.HandleMail(data); err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}
				}

				success = true
				break
			}

			if !success {
				fmt.Println("failed to process campaign event", zap.Error(lastError))
				continue
			}

			fmt.Println("campaign event processed successfully", zap.Any("event", event))
		}
	}
}

func ConsumerNotificationForSMS(brokers []string, config *sarama.Config) error {
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("failed to create consumer", zap.Error(err))
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(NotificationsTopic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("failed to create partition consumer", zap.Error(err))
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	const maxRetries = 3
	var event models.Notification

	for {
		select {
		case err := <-partitionConsumer.Errors():
			fmt.Println("error from partition consumer", zap.Error(err))
		case msg := <-partitionConsumer.Messages():
			var lastError error
			success := false

			for attempt := 0; attempt < maxRetries; attempt++ {
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				notificationData, err := json.Marshal(event)
				if err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				var notification models.Notification
				if err := json.Unmarshal(notificationData, &notification); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				if notification.Type == "sms" {
					notificationData, err := json.Marshal(notification.Data)
					if err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}

					var data map[string]interface{}
					if err := json.Unmarshal(notificationData, &data); err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}

					if err := sms.HandleSMS(data); err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}
				}

				success = true
				break
			}

			if !success {
				fmt.Println("failed to process campaign event", zap.Error(lastError))
				continue
			}

			fmt.Println("campaign event processed successfully", zap.Any("event", event))
		}
	}
}

func ConsumerNotificationForPush(brokers []string, config *sarama.Config) error {
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("failed to create consumer", zap.Error(err))
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(NotificationsTopic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("failed to create partition consumer", zap.Error(err))
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	const maxRetries = 3
	var event models.Notification

	for {
		select {
		case err := <-partitionConsumer.Errors():
			fmt.Println("error from partition consumer", zap.Error(err))
		case msg := <-partitionConsumer.Messages():
			var lastError error
			success := false

			for attempt := 0; attempt < maxRetries; attempt++ {
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				notificationData, err := json.Marshal(event)
				if err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				var notification models.Notification
				if err := json.Unmarshal(notificationData, &notification); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				if notification.Type == "push" {
					notificationData, err := json.Marshal(notification.Data)
					if err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}

					var data map[string]interface{}
					if err := json.Unmarshal(notificationData, &data); err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}

					if err := push.HandlePush(data); err != nil {
						lastError = err
						time.Sleep(time.Duration(attempt+1) * time.Second)
						continue
					}
				}

				success = true
				break
			}

			if !success {
				fmt.Println("failed to process campaign event", zap.Error(lastError))
				continue
			}

			fmt.Println("campaign event processed successfully", zap.Any("event", event))
		}
	}
}
