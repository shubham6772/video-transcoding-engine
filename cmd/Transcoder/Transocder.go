package main

import (
	"log"
	"video-processor/internal/config"
	"video-processor/internal/logger"
	"video-processor/internal/rabbitmq"
	"github.com/joho/godotenv"
)

type rmqConnection struct {
	rabbitURL string
	queueName string
}

type rmqEnvObject struct{
	rabbitURL string
	queueName string
}
func main() {

	godotenv.Load()

	rmqConfig := config.LoadRMQConfigObject()

	consumer, err := rabbitmq.NewConsumer(rmqConfig.RmqURL, rmqConfig.RmqQueue)
	if err != nil {
		logger.Error("Error Creating Consumer.. %v", err)
	}
	defer consumer.Close()

	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Waiting for messages...")

	for msg := range msgs {
		logger.Info("Received: %v", string(msg.Body))


		// your processing logic here
		// e.g. FFmpeg / HLS pipeline

		msg.Ack(false) //for single message acknowledgement
	}
}

