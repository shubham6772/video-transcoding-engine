package ffmpeg

import (
	"testing"
	"video-processor/internal/logger"
)

func Test_commandBuilder(t *testing.T) {
	
	resoltion := 1080
	videoPath := "/Users/shubham/Personal/video-processor/input.mp4"
	userIdHash := "abcdUser1"
	command, outputPath, err := CommandBuilder(resoltion, videoPath, userIdHash)
	
	if err != nil{
		logger.Error("%v", err);
	}
	logger.Info("command: %s, \n output path: %s", command, outputPath)

	ExecuteCommand(command);

}

func Test_checkResolutions(t *testing.T) {
	videoPath := "/Users/shubham/Personal/video-processor/input.mp4"

	got := CheckResolutions(videoPath)

	want := 1080

	if got != want {
		logger.Error("checkResolutions() = %v, want %v", got, want)
	}
}

func Test_executeCommand(t *testing.T){
	command := `echo "Hello world"`
	err := ExecuteCommand(command);

	if err != nil{
		logger.Error("Error: %v", err)
	}
}
