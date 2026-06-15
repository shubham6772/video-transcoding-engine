package config

import (
	"fmt"
	"os"
	"strconv"
	"video-processor/internal/constants"
)

type globalConfig struct {
	LogFilePath string
	VideoDir    string
	OutputDir   string
}

type loggerConfig struct {
	LogFilePath string
}

type videoCacheConfig struct {
	VideoFilePath string
}

type RmqConfigObject struct {
	RmqURL string
	RmqQueue string
}

type storageConfigObject struct {
	Endpoint string
	Accesskey string
	AccessSecret string
	BucketName string
	Secure bool
}

func getENV() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "local"
	}
	return env
}

func getValueFromEnv(envKey string) string {
	env := os.Getenv(envKey)
	if env == "" {
		errMsg := fmt.Sprintf("env not found for key: %v", envKey)
		panic(errMsg)
	}

	return env
}

func LoadLoggerConfig() *loggerConfig {
	env := getENV()
	switch env {
	case "stage":
		return &loggerConfig{
			LogFilePath: constants.LOGGING_BASE_PATH_VIDEO_PROCESSOR_STAGE,
		}
	case "prod":
		return &loggerConfig{
			LogFilePath: constants.LOGGING_BASE_PATH_VIDEO_PROCESSOR_PROD,
		}
	default:
		return &loggerConfig{
			LogFilePath: constants.LOGGING_BASE_PATH_VIDEO_PROCESSOR_LOCAL,
		}
	}
}

func LoadVideoCacheConfig() *videoCacheConfig{
	env := getENV()
	switch env {
	case "stage":
		return &videoCacheConfig{
			VideoFilePath: constants.VIDEO_BASE_PATH_FOLDER_STAGE,
		}
	case "prod":
		return &videoCacheConfig{
			VideoFilePath: constants.VIDEO_BASE_PATH_FOLDER_PROD,
		}
	default:
		return &videoCacheConfig{
			VideoFilePath: constants.VIDEO_BASE_PATH_FOLDER_LOCAL,
		}
	}

}

func LoadRMQConfigObject() *RmqConfigObject {
	
	return &RmqConfigObject{
		RmqURL: getValueFromEnv("RMQ_URL"),
		RmqQueue: getValueFromEnv("RMQ_QUEUE"),
	}

}

func LoadStorageConfigObject() *storageConfigObject{
	val := getValueFromEnv("MINIO_USE_SSL")
	isSecure, err := strconv.ParseBool(val)
	if err != nil{
		panic(err)
	}
	return &storageConfigObject{
		Endpoint : getValueFromEnv("MINIO_ENDPOINT"),
		Accesskey: getValueFromEnv("MINIO_ACCESS_KEY"),
		AccessSecret: getValueFromEnv("MINIO_SECRET_KEY"),
		BucketName: getValueFromEnv("MINIO_BUCKET"),
		Secure: isSecure,
	}
}
