package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"video-processor/internal/config"
	"video-processor/internal/ffmpeg"
	"video-processor/internal/logger"
	"video-processor/internal/rabbitmq"
	"video-processor/internal/storage"
)

const TRANSCODER_VERSION = "1.0.0"

type VideoJob struct {
	VideoID   string `json:"videoId"`
	FilePath  string `json:"filePath"`
	UserID    string `json:"userId"`
	Extension string `json:"ext"`
}

func main() {
	_ = godotenv.Load()

	rmqConfig := config.LoadRMQConfigObject()

	consumer, err := rabbitmq.NewConsumer(
		rmqConfig.RmqURL,
		rmqConfig.RmqQueue,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	storageConfig := config.LoadStorageConfigObject()

	minioClient, storageErr := minio.New(storageConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(storageConfig.Accesskey, storageConfig.AccessSecret, ""),
		Secure: storageConfig.Secure,
	})

	if storageErr != nil {
		log.Fatalf("%v", storageErr.Error())
	}

	store := &storage.MinIOStorage{
		Client: minioClient,
	}

	manager := &storage.Manager{
		Store: store,
	}

	cacheStoragePath := config.LoadVideoCacheConfig()

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
		logger.Info("Extension: %s", job.Extension)

		//donwload file locally

		baseDir := filepath.Join(cacheStoragePath.VideoFilePath, job.UserID, job.VideoID)
		inputPath := filepath.Join(baseDir, "input", "source."+job.Extension)
		outputPath := filepath.Join(baseDir, "output/")

		if err := os.MkdirAll(outputPath, 0755); err != nil {
			msg.Ack(false)
			continue
		}

		dnldErr := manager.DownloadFile(storageConfig.BucketName, job.FilePath, inputPath)

		if dnldErr != nil {
			logger.Error("download error occured: %v", dnldErr)
			msg.Ack(false)
			continue
		}

		logger.Info("checking resolution of video %v", job.VideoID)
		currentResolution, reserr := ffmpeg.CheckResolution(inputPath)

		if reserr != nil {
			logger.Error("unsupported resolution of video file: %s", job.VideoID)
			msg.Ack(false)
			continue
		}

		transcode_command, output_path, cmderr := ffmpeg.CommandBuilder(currentResolution, inputPath, outputPath)

		if cmderr != nil {
			logger.Error("Failed to generate transcode command for video: %v, userid: %v, path: %v err: %v", job.VideoID, job.UserID, inputPath, err)
			msg.Ack(false)
			continue
		}

		logger.Info("transcode command: %v", transcode_command)
		logger.Info("executing command:")
		execerr := ffmpeg.ExecuteCommand(transcode_command)

		if execerr != nil {
			logger.Error("failed to execute: %v", transcode_command)
			msg.Ack(false)
			continue
		}

		err = manager.UploadFolder(
			storageConfig.BucketName,
			filepath.Join(job.UserID, job.VideoID, "output"),
			outputPath,
		)

		logger.Info("video file generated at location: %s", output_path)

		err = msg.Ack(false)
		dltErr := DeleteLocalFolder(filepath.Join(cacheStoragePath.VideoFilePath, job.UserID))
		if dltErr != nil {
			continue
		}

		if err != nil {
			logger.Error("ack failed: %v", err)
			continue
		}

		logger.Info("Waiting for messages...")

	}
}

func DeleteLocalFolder(path string) error {
	return os.RemoveAll(path)
}