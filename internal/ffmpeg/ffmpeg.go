package ffmpeg

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"video-processor/internal/config"
	"video-processor/internal/logger"
)

type Rendition struct {
	Name       string
	Width      int
	Height     int
	Bitrate    string
	MaxRate    string
	BufferSize string
}

var renditions = []Rendition{
	{
		Name:       "1080p",
		Width:      1920,
		Height:     1080,
		Bitrate:    "5000k",
		MaxRate:    "5350k",
		BufferSize: "7500k",
	},
	{
		Name:       "720p",
		Width:      1280,
		Height:     720,
		Bitrate:    "2800k",
		MaxRate:    "2996k",
		BufferSize: "4200k",
	},
	{
		Name:       "480p",
		Width:      854,
		Height:     480,
		Bitrate:    "1400k",
		MaxRate:    "1498k",
		BufferSize: "2100k",
	},
	{
		Name:       "360p",
		Width:      640,
		Height:     360,
		Bitrate:    "800k",
		MaxRate:    "856k",
		BufferSize: "1200k",
	},
	{
		Name:       "240p",
		Width:      426,
		Height:     240,
		Bitrate:    "400k",
		MaxRate:    "428k",
		BufferSize: "600k",
	},
}

func CommandBuilder(maxResolution int, videoPath string, videoIdHash string) (string, string, error) {
	videoConfig := config.LoadVideoCacheConfig()

	outputPath := filepath.Join(
		videoConfig.VideoFilePath,
		videoIdHash,
	)

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return "", "", fmt.Errorf("create output directory: %w", err)
	}

	var selected []Rendition

	found := false
	for _, r := range renditions {
		if r.Height == maxResolution {
			found = true
		}

		if found {
			selected = append(selected, r)
		}
	}

	if len(selected) == 0 {
		return "", "", fmt.Errorf("unsupported resolution: %d", maxResolution)
	}

	var splitOutputs []string

	for _, r := range selected {
		splitOutputs = append(
			splitOutputs,
			fmt.Sprintf("[v%d]", r.Height),
		)
	}

	filter := fmt.Sprintf(
		"[0:v]split=%d%s;",
		len(selected),
		strings.Join(splitOutputs, ""),
	)

	for _, r := range selected {
		filter += fmt.Sprintf(
			"[v%d]scale=%d:%d[v%do];",
			r.Height,
			r.Width,
			r.Height,
			r.Height,
		)
	}

	var maps []string

	for idx, r := range selected {
		maps = append(
			maps,
			fmt.Sprintf(
				`-map "[v%[1]do]" -map a:0 `+
					`-c:v:%[2]d libx264 `+
					`-b:v:%[2]d %[3]s `+
					`-maxrate:v:%[2]d %[4]s `+
					`-bufsize:v:%[2]d %[5]s`,
				r.Height,
				idx,
				r.Bitrate,
				r.MaxRate,
				r.BufferSize,
			),
		)
	}

	var streamMap []string

	for idx, r := range selected {
		streamMap = append(
			streamMap,
			fmt.Sprintf(
				"v:%d,a:%d,name:%s",
				idx,
				idx,
				r.Name,
			),
		)
	}

	cmd := fmt.Sprintf(
		`ffmpeg -i "%s" \
		-filter_complex "%s" \
		%s \
		-c:a aac \
		-b:a 128k \
		-g 48 \
		-keyint_min 48 \
		-sc_threshold 0 \
		-hls_time 6 \
		-hls_playlist_type vod \
		-hls_flags independent_segments \
		-master_pl_name master.m3u8 \
		-var_stream_map "%s" \
		-hls_segment_filename "%s/%%v/seg_%%03d.ts" \
		-f hls \
		-progress pipe:1 \
		-stats_period 1 \
		"%s/%%v/index.m3u8"`,
		videoPath,
		filter,
		strings.Join(maps, " "),
		strings.Join(streamMap, " "),
		outputPath,
		outputPath,
	)

	return cmd, outputPath, nil
}

func CheckResolution(videoPath string) (int, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=height",
		"-of", "csv=p=0",
		videoPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf(
			"ffprobe failed for %s: %w, output: %s",
			videoPath,
			err,
			string(output),
		)
	}

	heightStr := strings.TrimSpace(string(output))

	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to parse height '%s': %w",
			heightStr,
			err,
		)
	}

	logger.Info("video resolution: %d", height)

	return height, nil
}


func ExecuteCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			logger.Info("[FFMPEG] %s", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			logger.Error("stdout scanner error: %v", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)

		for scanner.Scan() {
			logger.Info("[FFMPEG] %s", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			logger.Error("stderr scanner error: %v", err)
		}
	}()

	if err := cmd.Wait(); err != nil {
		logger.Error("ffmpeg failed: %v", err)
		return err
	}

	logger.Info("ffmpeg completed successfully")

	return nil
}