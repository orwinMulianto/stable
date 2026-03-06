package utils

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"time"
)

var (
	Log       *slog.Logger
	ErrorLog  *slog.Logger
	Important *slog.Logger
)

// LogEntry represents a standardized log entry for analytics and ML
type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Domain    string         `json:"domain"`
	Action    string         `json:"action"`
	Message   string         `json:"message"`
	RequestID string         `json:"request_id,omitempty"`
	UserID    uint           `json:"user_id,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

func openLogFile(path string) *os.File {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	return file
}

func InitLogger() {
	_ = os.MkdirAll("logs", 0755)

	// files
	appFile := openLogFile("logs/app.log")
	errorFile := openLogFile("logs/error.log")
	importantFile := openLogFile("logs/important.log")

	// stdout + file
	appWriter := io.MultiWriter(os.Stdout, appFile)

	Log = slog.New(
		slog.NewJSONHandler(appWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	ErrorLog = slog.New(
		slog.NewJSONHandler(errorFile, &slog.HandlerOptions{
			Level: slog.LevelError,
		}),
	)

	Important = slog.New(
		slog.NewJSONHandler(importantFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	Log.Info("Logger initialized")
}

// writeLog writes a structured log entry to the appropriate log file
func writeLog(entry LogEntry) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		Log.Error("Failed to marshal log entry", "error", err.Error())
		return
	}

	switch entry.Level {
	case "error":
		ErrorLog.Info(string(jsonData))
		Log.Error(entry.Message,
			"domain", entry.Domain,
			"action", entry.Action,
			"request_id", entry.RequestID,
			"user_id", entry.UserID,
		)
	case "success", "info":
		Log.Info(entry.Message,
			"domain", entry.Domain,
			"action", entry.Action,
			"request_id", entry.RequestID,
			"user_id", entry.UserID,
		)
	case "warning":
		Log.Warn(entry.Message,
			"domain", entry.Domain,
			"action", entry.Action,
			"request_id", entry.RequestID,
			"user_id", entry.UserID,
		)
	}

	// Write structured entry to important.log for analytics/ML
	Important.Info(string(jsonData))
}

// LogSuccess logs a successful operation
func LogSuccess(domain, action, message, requestID string, userID uint, metadata map[string]any) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "success",
		Domain:    domain,
		Action:    action,
		Message:   message,
		RequestID: requestID,
		UserID:    userID,
		Metadata:  metadata,
	}
	writeLog(entry)
}

// LogInfo logs an informational message
func LogInfo(domain, action, message, requestID string, userID uint, metadata map[string]any) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "info",
		Domain:    domain,
		Action:    action,
		Message:   message,
		RequestID: requestID,
		UserID:    userID,
		Metadata:  metadata,
	}
	writeLog(entry)
}

// LogError logs an error
func LogError(domain, action, message, requestID string, userID uint, metadata map[string]any) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "error",
		Domain:    domain,
		Action:    action,
		Message:   message,
		RequestID: requestID,
		UserID:    userID,
		Metadata:  metadata,
	}
	writeLog(entry)
}

// LogWarning logs a warning message
func LogWarning(domain, action, message, requestID string, userID uint, metadata map[string]any) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "warning",
		Domain:    domain,
		Action:    action,
		Message:   message,
		RequestID: requestID,
		UserID:    userID,
		Metadata:  metadata,
	}
	writeLog(entry)
}
