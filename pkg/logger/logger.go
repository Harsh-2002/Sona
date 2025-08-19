package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	logFile *os.File
	logger  *log.Logger
)

// InitLogger initializes the logger with a file in .sona folder
func InitLogger() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	sonaDir := filepath.Join(homeDir, ".sona")
	if err := os.MkdirAll(sonaDir, 0755); err != nil {
		return fmt.Errorf("failed to create .sona directory: %v", err)
	}

	logPath := filepath.Join(sonaDir, "sona.log")
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	logger = log.New(logFile, "", log.LstdFlags)
	return nil
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// LogInfo logs an info message
func LogInfo(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf("[INFO] "+format, args...)
	}
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf("[ERROR] "+format, args...)
	}
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf("[DEBUG] "+format, args...)
	}
}

// LogWarning logs a warning message
func LogWarning(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf("[WARNING] "+format, args...)
	}
}

// GetLogPath returns the path to the log file
func GetLogPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".sona", "sona.log")
}

// LogCommand logs a command execution
func LogCommand(cmd string, args []string, output string, err error) {
	if logger != nil {
		logger.Printf("[COMMAND] %s %v", cmd, args)
		if output != "" {
			logger.Printf("[OUTPUT] %s", output)
		}
		if err != nil {
			logger.Printf("[ERROR] %v", err)
		}
	}
}
