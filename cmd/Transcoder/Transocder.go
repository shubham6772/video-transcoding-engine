package main

import (
	"encoding/json"
	"log"

	"github.com/joho/godotenv"

	"video-processor/internal/config"
	"video-processor/internal/ffmpeg"
	"video-processor/internal/logger"
	"video-processor/internal/rabbitmq"
)

const TRANSCODER_VERSION = "1.0.0"

type VideoJob struct {
	VideoID  string `json:"videoId"`
	FilePath string `json:"filePath"`
	UserID   string `json:"userId"`
}

func main() {
	godotenv.Load()

	rmqConfig := config.LoadRMQConfigObject()

	consumer, err := rabbitmq.NewConsumer(
		rmqConfig.RmqURL,
		rmqConfig.RmqQueue,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Waiting for messages...")

	for msg := range msgs {
		var job VideoJob

		err := json.Unmarshal(msg.Body, &job)
		if err != nil {
			logger.Error("Failed to parse message: %v", err)

			// Invalid message, discard it
			msg.Nack(false, false)

			continue
		}

		logger.Info("Received Video Job")
		logger.Info("VideoID: %s", job.VideoID)
		logger.Info("FilePath: %s", job.FilePath)
		logger.Info("UserID: %s", job.UserID)


		logger.Info("checking resolution of video %v", job.VideoID)
		currentResolution, reserr := ffmpeg.CheckResolution(job.FilePath)

		if reserr != nil{
			logger.Error("unsupported resolution of video file: %s", job.VideoID)
			msg.Ack(false)
			continue
		}

		transcode_command, output_path, cmderr := ffmpeg.CommandBuilder(currentResolution, job.FilePath, job.UserID)

		if cmderr != nil{
			logger.Error("Failed to generate transcode command for video: %v, userid: %v, path: %v err: %v", job.VideoID, job.UserID, job.FilePath, err)
			msg.Ack(false)
			continue 
		}

		logger.Info("transcode command: %v", transcode_command);
		logger.Info("executing command:")
		execerr := ffmpeg.ExecuteCommand(transcode_command)
		
		if execerr != nil{
			logger.Error("failed to execute: %v", transcode_command)
			msg.Ack(false)
			continue
		}
		logger.Info("video file generated at location: %s", output_path)

		err = msg.Ack(false)
		if err != nil {
			logger.Error("ack failed: %v", err)
			return
		}

		logger.Info("Waiting for messages...")
		
	}
}