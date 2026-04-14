package logger

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// Logger is a wrapper around slog.Logger with component tracking
type Logger struct {
	logger    *slog.Logger
	component string
}

// LogEntry represents a single log entry with all metadata
type LogEntry struct {
	Level      LogLevel      `json:"level"`
	Timestamp  string        `json:"timestamp"`
	Component  string        `json:"component"`
	Message    string        `json:"message"`
	Error      string        `json:"error,omitempty"`
	StackTrace string        `json:"stack_trace,omitempty"`
	Attributes LogAttributes `json:"attributes,omitempty"`
}

// LogAttributes contains structured logging attributes
type LogAttributes struct {
	UserID   string `json:"user_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	Method   string `json:"method,omitempty"`
	Status   int    `json:"status,omitempty"`
	Reason   string `json:"reason,omitempty"`
	Details  string `json:"details,omitempty"`
}

// AutoFlushWriter writes and immediately flushes to ensure real-time output
type AutoFlushWriter struct {
	writer      io.Writer
	flusher     *bufio.Writer
	isBuffered  bool
}

// Write writes data and flushes immediately if buffered
func (w *AutoFlushWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if err != nil {
		return n, err
	}
	if w.isBuffered && w.flusher != nil {
		_ = w.flusher.Flush()
	}
	return n, nil
}

var (
	// Global logger instance
	globalLogger *Logger
	// Log file for persistence
	logFile *os.File
	// Buffered writer for file to ensure immediate flushing
	bufferedLogFile *bufio.Writer
)

// InitLogger initializes the global logger with file and console output
func InitLogger(logPath string) error {
	// Create logs directory if it doesn't exist
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open or create log file with append mode
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logFile = file

	// Create multi-writer: write to both file and console
	// Use auto-flush writer for file to ensure real-time output
	bufferedLogFile = bufio.NewWriter(logFile)
	fileWriter := &AutoFlushWriter{
		writer:     bufferedLogFile,
		flusher:    bufferedLogFile,
		isBuffered: true,
	}

	// Create multi-writer with auto-flushing file and stdout
	multiWriter := io.MultiWriter(fileWriter, os.Stdout)

	// Create JSON handler for structured logging
	jsonHandler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			// Customize attribute formatting
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(time.Now().UTC().Format(time.RFC3339)),
				}
			}
			return a
		},
	})

	logger := slog.New(jsonHandler)
	globalLogger = &Logger{
		logger:    logger,
		component: "TravelLink",
	}

	return nil
}

// GetLogger returns a logger instance with specified component
func GetLogger(component string) *Logger {
	if globalLogger == nil {
		// Fallback to console-only logger if not initialized
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		globalLogger = &Logger{
			logger:    slog.New(handler),
			component: "TravelLink",
		}
	}

	return &Logger{
		logger:    globalLogger.logger,
		component: component,
	}
}

// CloseLogger closes the log file
func CloseLogger() error {
	if bufferedLogFile != nil {
		if err := bufferedLogFile.Flush(); err != nil {
			return err
		}
	}
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// captureStackTrace captures stack trace from the current call
func captureStackTrace(skip int) string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// Error logs an error with full stack trace
func (l *Logger) Error(ctx context.Context, message string, err error, attrs LogAttributes) {
	entry := LogEntry{
		Level:      ERROR,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Component:  l.component,
		Message:    message,
		StackTrace: captureStackTrace(2),
		Attributes: attrs,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	l.logger.ErrorContext(ctx, message,
		slog.String("component", l.component),
		slog.String("error", entry.Error),
		slog.String("stack_trace", entry.StackTrace),
		slog.String("user_id", attrs.UserID),
		slog.String("email", attrs.Email),
		slog.String("endpoint", attrs.Endpoint),
		slog.String("method", attrs.Method),
		slog.Int("status", attrs.Status),
		slog.String("reason", attrs.Reason),
		slog.String("details", attrs.Details),
	)
}

// AuthorizationDenied logs when an authorization attempt is denied
func (l *Logger) AuthorizationDenied(ctx context.Context, userID, email, endpoint, method, reason string) {
	l.logger.ErrorContext(ctx, "Authorization Denied",
		slog.String("component", l.component),
		slog.String("user_id", userID),
		slog.String("email", email),
		slog.String("endpoint", endpoint),
		slog.String("method", method),
		slog.String("reason", reason),
		slog.String("type", "AUTHORIZATION_DENIED"),
	)
}

// PermissionDenied logs when a permission check fails
func (l *Logger) PermissionDenied(ctx context.Context, userID, email, resource, action, reason string) {
	l.logger.ErrorContext(ctx, "Permission Denied",
		slog.String("component", l.component),
		slog.String("user_id", userID),
		slog.String("email", email),
		slog.String("resource", resource),
		slog.String("action", action),
		slog.String("reason", reason),
		slog.String("type", "PERMISSION_DENIED"),
	)
}

// AccessAttemptDenied logs failed access attempts (invalid credentials, etc)
func (l *Logger) AccessAttemptDenied(ctx context.Context, email, endpoint, method, reason string) {
	l.logger.ErrorContext(ctx, "Access Attempt Denied",
		slog.String("component", l.component),
		slog.String("email", email),
		slog.String("endpoint", endpoint),
		slog.String("method", method),
		slog.String("reason", reason),
		slog.String("type", "ACCESS_DENIED"),
	)
}

// ValidationFailed logs validation failures
func (l *Logger) ValidationFailed(ctx context.Context, field, reason string, attrs LogAttributes) {
	l.logger.ErrorContext(ctx, "Validation Failed",
		slog.String("component", l.component),
		slog.String("user_id", attrs.UserID),
		slog.String("email", attrs.Email),
		slog.String("field", field),
		slog.String("reason", reason),
		slog.String("type", "VALIDATION_ERROR"),
	)
}

// Info logs informational messages
func (l *Logger) Info(ctx context.Context, message string, attrs LogAttributes) {
	l.logger.InfoContext(ctx, message,
		slog.String("component", l.component),
		slog.String("user_id", attrs.UserID),
		slog.String("email", attrs.Email),
		slog.String("endpoint", attrs.Endpoint),
		slog.String("method", attrs.Method),
		slog.Int("status", attrs.Status),
		slog.String("details", attrs.Details),
	)
}

// Debug logs debug-level messages
func (l *Logger) Debug(ctx context.Context, message string, attrs LogAttributes) {
	l.logger.DebugContext(ctx, message,
		slog.String("component", l.component),
		slog.String("user_id", attrs.UserID),
		slog.String("details", attrs.Details),
	)
}

// Warn logs warning-level messages
func (l *Logger) Warn(ctx context.Context, message string, attrs LogAttributes) {
	l.logger.WarnContext(ctx, message,
		slog.String("component", l.component),
		slog.String("user_id", attrs.UserID),
		slog.String("reason", attrs.Reason),
		slog.String("details", attrs.Details),
	)
}

// DatabaseError logs database operation errors
func (l *Logger) DatabaseError(ctx context.Context, operation string, err error, userID string) {
	l.logger.ErrorContext(ctx, fmt.Sprintf("Database Error: %s", operation),
		slog.String("component", l.component),
		slog.String("operation", operation),
		slog.String("error", err.Error()),
		slog.String("stack_trace", captureStackTrace(2)),
		slog.String("user_id", userID),
		slog.String("type", "DATABASE_ERROR"),
	)
}

// AuthenticationError logs authentication-related errors
func (l *Logger) AuthenticationError(ctx context.Context, email, reason string) {
	l.logger.ErrorContext(ctx, "Authentication Error",
		slog.String("component", l.component),
		slog.String("email", email),
		slog.String("reason", reason),
		slog.String("type", "AUTHENTICATION_ERROR"),
	)
}

// RequestLogging logs incoming request details
func (l *Logger) RequestLogging(ctx context.Context, method, endpoint, userID string, status int) {
	l.logger.InfoContext(ctx, "Request",
		slog.String("component", l.component),
		slog.String("method", method),
		slog.String("endpoint", endpoint),
		slog.String("user_id", userID),
		slog.Int("status", status),
		slog.String("type", "REQUEST"),
	)
}

// SecurityEvent logs security-related events
func (l *Logger) SecurityEvent(ctx context.Context, eventType, userID, email, details string) {
	l.logger.WarnContext(ctx, fmt.Sprintf("Security Event: %s", eventType),
		slog.String("component", l.component),
		slog.String("event_type", eventType),
		slog.String("user_id", userID),
		slog.String("email", email),
		slog.String("details", details),
		slog.String("type", "SECURITY_EVENT"),
	)
}

// ServiceError logs service layer errors
func (l *Logger) ServiceError(ctx context.Context, service, operation string, err error, userID string) {
	l.logger.ErrorContext(ctx, fmt.Sprintf("Service Error: %s.%s", service, operation),
		slog.String("component", l.component),
		slog.String("service", service),
		slog.String("operation", operation),
		slog.String("error", err.Error()),
		slog.String("stack_trace", captureStackTrace(2)),
		slog.String("user_id", userID),
		slog.String("type", "SERVICE_ERROR"),
	)
}

// UserRegistration logs when a new user registers
func (l *Logger) UserRegistration(ctx context.Context, userID, email, username string) {
	l.logger.InfoContext(ctx, "User Registration",
		slog.String("component", l.component),
		slog.String("user_id", userID),
		slog.String("email", email),
		slog.String("username", username),
		slog.String("type", "USER_REGISTRATION"),
	)
}

// UserAuthentication logs when a user successfully authenticates
func (l *Logger) UserAuthentication(ctx context.Context, userID, email string) {
	l.logger.InfoContext(ctx, "User Authentication",
		slog.String("component", l.component),
		slog.String("user_id", userID),
		slog.String("email", email),
		slog.String("type", "USER_AUTHENTICATION"),
	)
}

// UserLogout logs when a user logs out
func (l *Logger) UserLogout(ctx context.Context, userID string) {
	l.logger.InfoContext(ctx, "User Logout",
		slog.String("component", l.component),
		slog.String("user_id", userID),
		slog.String("type", "USER_LOGOUT"),
	)
}

// AuthorizationError logs when an authorization error occurs
func (l *Logger) AuthorizationError(ctx context.Context, userID, reason string) {
	l.logger.ErrorContext(ctx, "Authorization Error",
		slog.String("component", l.component),
		slog.String("user_id", userID),
		slog.String("reason", reason),
		slog.String("type", "AUTHORIZATION_ERROR"),
	)
}

// ValidationError logs when validation fails with field-level errors
func (l *Logger) ValidationError(ctx context.Context, errors map[string]string) {
	errorStr := ""
	for field, err := range errors {
		if errorStr != "" {
			errorStr += "; "
		}
		errorStr += fmt.Sprintf("%s: %s", field, err)
	}
	l.logger.ErrorContext(ctx, "Validation Error",
		slog.String("component", l.component),
		slog.String("errors", errorStr),
		slog.String("type", "VALIDATION_ERROR"),
	)
}
