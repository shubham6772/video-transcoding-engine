package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"video-processor/internal/config"
)

var (
	debugLogger *log.Logger
	errorLogger *log.Logger
	infoLogger  *log.Logger

	debugFile *os.File
	errorFile *os.File

	currentDay string

	mu   sync.Mutex
	once sync.Once
	err  error
)

func Init() error {
	var e error

	once.Do(func() {
		e = initLogger()
	})

	if e != nil {
		return e
	}

	// ensure correct day state (safe, not rotation hack)
	return ensureToday()
}

func initLogger() error {
	basePath := config.LoadLoggerConfig().LogFilePath

	now := time.Now()
	currentDay = now.Format("2006-01-02")

	dir := filepath.Join(
		basePath,
		now.Format("2006"),
		now.Format("01"),
		now.Format("02"),
	)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	var err error

	debugFile, err = os.OpenFile(
		filepath.Join(dir, "debug.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	errorFile, err = os.OpenFile(
		filepath.Join(dir, "error.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		_ = debugFile.Close()
		return err
	}

	debugWriter := io.MultiWriter(os.Stdout, debugFile)
	errorWriter := io.MultiWriter(os.Stderr, errorFile)

	infoLogger = log.New(debugWriter, "INFO: ", log.Ldate|log.Ltime)
	debugLogger = log.New(debugWriter, "DEBUG: ", log.Ldate|log.Ltime)
	errorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime)

	return nil
}

func ensureToday() error {
	today := time.Now().Format("2006-01-02")

	if today == currentDay {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	if today == currentDay {
		return nil
	}

	closeFiles()
	return initLogger()
}

func closeFiles() {
	if debugFile != nil {
		_ = debugFile.Close()
		debugFile = nil
	}

	if errorFile != nil {
		_ = errorFile.Close()
		errorFile = nil
	}
}

func Debug(format string, v ...any) {
	_ = ensureToday()

	_, file, line, ok := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)

	if !ok {
		debugLogger.Println(msg)
		return
	}

	log.Printf("%s:%d: %s", file, line, msg)
	debugLogger.Printf("%s:%d: %s", file, line, msg)
}

func Info(format string, v ...any) {
	_ = ensureToday()

	_, file, line, ok := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)

	if !ok {
		infoLogger.Println(msg)
		return
	}

	log.Printf("%s:%d: %s", file, line, msg)
	infoLogger.Printf("%s:%d: %s", file, line, msg)
}

func Error(format string, v ...any) {
	_ = ensureToday()

	_, file, line, ok := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)

	if !ok {
		errorLogger.Println(msg)
		return
	}

	log.Printf("%s:%d: %s", file, line, msg)
	errorLogger.Printf("%s:%d: %s", file, line, msg)
}

func Close() {
	mu.Lock()
	defer mu.Unlock()
	closeFiles()
}
