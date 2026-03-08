package logger

import (
	"discordcommandbot/pkg/errors"
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelCritical
)

// Logger is the application logger with level-based filtering
type Logger struct {
	debugMode   bool
	minLevel    Level
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var (
	// Default is the default logger instance
	Default *Logger
)

// Init initializes the default logger
func Init(debugMode bool) {
	minLevel := LevelInfo
	if debugMode {
		minLevel = LevelDebug
	}

	Default = &Logger{
		debugMode:   debugMode,
		minLevel:    minLevel,
		infoLogger:  log.New(os.Stdout, "", 0),
		warnLogger:  log.New(os.Stdout, "", 0),
		errorLogger: log.New(os.Stderr, "", 0),
		debugLogger: log.New(os.Stdout, "", 0),
	}

	if debugMode {
		Default.Info("🔧 Debug mode enabled")
	}
}

// formatMessage formats a log message with timestamp and level
func formatMessage(level, emoji, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] %s [%s] %s", timestamp, emoji, level, message)
}

// ============================================================================
// DEBUG LEVEL (only shown when DEBUG_MODE=true)
// ============================================================================

// Debug logs a debug message (only if debug mode is enabled)
func (l *Logger) Debug(format string, args ...any) {
	if l.debugMode && l.minLevel <= LevelDebug {
		message := fmt.Sprintf(format, args...)
		l.debugLogger.Println(formatMessage("DEBUG", "🔍", message))
	}
}

// DebugError logs an Error at debug level
func (l *Logger) DebugError(err *errors.Error) {
	if l.debugMode && l.minLevel <= LevelDebug {
		message := err.Error()
		if len(err.Context) > 0 {
			message += fmt.Sprintf(" | Context: %v", err.Context)
		}
		l.debugLogger.Println(formatMessage("DEBUG", "🔍", message))
	}
}

// ============================================================================
// INFO LEVEL
// ============================================================================

// Info logs an informational message
func (l *Logger) Info(format string, args ...any) {
	if l.minLevel <= LevelInfo {
		message := fmt.Sprintf(format, args...)
		l.infoLogger.Println(formatMessage("INFO", "✅", message))
	}
}

// InfoError logs an Error at info level
func (l *Logger) InfoError(err *errors.Error) {
	if l.minLevel <= LevelInfo {
		message := err.Error()
		if len(err.Context) > 0 {
			message += fmt.Sprintf(" | Context: %v", err.Context)
		}
		l.infoLogger.Println(formatMessage("INFO", "ℹ️", message))
	}
}

// ============================================================================
// WARN LEVEL
// ============================================================================

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...any) {
	if l.minLevel <= LevelWarn {
		message := fmt.Sprintf(format, args...)
		l.warnLogger.Println(formatMessage("WARN", "⚠️", message))
	}
}

// WarnError logs an Error at warn level
func (l *Logger) WarnError(err *errors.Error) {
	if l.minLevel <= LevelWarn {
		message := err.Error()
		if len(err.Context) > 0 {
			message += fmt.Sprintf(" | Context: %v", err.Context)
		}
		l.warnLogger.Println(formatMessage("WARN", "⚠️", message))
	}
}

// ============================================================================
// CRITICAL LEVEL
// ============================================================================

// Critical logs a critical error message
func (l *Logger) Critical(format string, args ...any) {
	if l.minLevel <= LevelCritical {
		message := fmt.Sprintf(format, args...)
		l.errorLogger.Println(formatMessage("CRITICAL", "❌", message))
	}
}

// CriticalError logs an Error at critical level
func (l *Logger) CriticalError(err *errors.Error) {
	if l.minLevel <= LevelCritical {
		message := err.Error()
		if len(err.Context) > 0 {
			message += fmt.Sprintf(" | Context: %v", err.Context)
		}
		l.errorLogger.Println(formatMessage("CRITICAL", "❌", message))
	}
}

// ============================================================================
// AUTO-DETECT ERROR LEVEL (logs based on Error.Level)
// ============================================================================

// LogError logs an Error at its appropriate level
func (l *Logger) LogError(err *errors.Error) {
	switch err.Level {
	case errors.DebugLevel:
		l.DebugError(err)
	case errors.InfoLevel:
		l.InfoError(err)
	case errors.WarnLevel:
		l.WarnError(err)
	case errors.CriticalLevel:
		l.CriticalError(err)
	default:
		l.WarnError(err)
	}
}

// ============================================================================
// GLOBAL CONVENIENCE FUNCTIONS
// ============================================================================

// Debug logs a debug message using the default logger
func Debug(format string, args ...any) {
	if Default != nil {
		Default.Debug(format, args...)
	}
}

// Info logs an info message using the default logger
func Info(format string, args ...any) {
	if Default != nil {
		Default.Info(format, args...)
	}
}

// Warn logs a warn message using the default logger
func Warn(format string, args ...any) {
	if Default != nil {
		Default.Warn(format, args...)
	}
}

// Critical logs a critical message using the default logger
func Critical(format string, args ...any) {
	if Default != nil {
		Default.Critical(format, args...)
	}
}

// LogError logs an Error using the default logger
func LogError(err *errors.Error) {
	if Default != nil {
		Default.LogError(err)
	}
}
