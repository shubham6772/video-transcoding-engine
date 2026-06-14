package config

import (
	"fmt"
	"os"
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
