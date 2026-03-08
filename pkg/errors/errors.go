package errors

import (
	"fmt"
	"time"
)

// Level represents the severity of an error
type Level string

const (
	DebugLevel    Level = "DEBUG"
	InfoLevel     Level = "INFO"
	WarnLevel     Level = "WARN"
	CriticalLevel Level = "CRITICAL"
)

// Category represents the category/type of an error
type Category string

const (
	Validation  Category = "VALIDATION"
	Config      Category = "CONFIG"
	Integration Category = "INTEGRATION"
	Service     Category = "SERVICE"
	Handler     Category = "HANDLER"
	Database    Category = "DATABASE"
	API         Category = "API"
	Network     Category = "NETWORK"
	Auth        Category = "AUTH"
	Unknown     Category = "UNKNOWN"
)

// Error represents a custom application error with context
type Error struct {
	Level     Level
	Category  Category
	Message   string
	Cause     error
	Timestamp time.Time
	Context   map[string]any
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s][%s] %s: %v", e.Level, e.Category, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s][%s] %s", e.Level, e.Category, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Cause
}

// WithContext adds context to the error
func (e *Error) WithContext(key string, value any) *Error {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// New creates a new Error with the specified parameters
func New(level Level, category Category, message string, cause error) *Error {
	return &Error{
		Level:     level,
		Category:  category,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
		Context:   make(map[string]any),
	}
}

// ============================================================================
// VALIDATION ERRORS
// ============================================================================

// NewValidation creates a new validation error (INFO level)
func NewValidation(message string, cause error) *Error {
	return New(InfoLevel, Validation, message, cause)
}

// NewValidationf creates a new validation error with formatted message
func NewValidationf(format string, args ...any) *Error {
	return New(InfoLevel, Validation, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// CONFIG ERRORS
// ============================================================================

// NewConfig creates a new configuration error (CRITICAL level)
func NewConfig(message string, cause error) *Error {
	return New(CriticalLevel, Config, message, cause)
}

// NewConfigf creates a new configuration error with formatted message
func NewConfigf(format string, args ...any) *Error {
	return New(CriticalLevel, Config, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// INTEGRATION ERRORS
// ============================================================================

// NewIntegration creates a new integration error (WARN level)
func NewIntegration(message string, cause error) *Error {
	return New(WarnLevel, Integration, message, cause)
}

// NewIntegrationf creates a new integration error with formatted message
func NewIntegrationf(format string, args ...any) *Error {
	return New(WarnLevel, Integration, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// SERVICE ERRORS
// ============================================================================

// NewService creates a new service error (WARN level)
func NewService(message string, cause error) *Error {
	return New(WarnLevel, Service, message, cause)
}

// NewServicef creates a new service error with formatted message
func NewServicef(format string, args ...any) *Error {
	return New(WarnLevel, Service, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// HANDLER ERRORS
// ============================================================================

// NewHandler creates a new handler error (WARN level)
func NewHandler(message string, cause error) *Error {
	return New(WarnLevel, Handler, message, cause)
}

// NewHandlerf creates a new handler error with formatted message
func NewHandlerf(format string, args ...any) *Error {
	return New(WarnLevel, Handler, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// DATABASE ERRORS
// ============================================================================

// NewDatabase creates a new database error (CRITICAL level)
func NewDatabase(message string, cause error) *Error {
	return New(CriticalLevel, Database, message, cause)
}

// NewDatabasef creates a new database error with formatted message
func NewDatabasef(format string, args ...any) *Error {
	return New(CriticalLevel, Database, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// API ERRORS
// ============================================================================

// NewAPI creates a new API error (WARN level)
func NewAPI(message string, cause error) *Error {
	return New(WarnLevel, API, message, cause)
}

// NewAPIf creates a new API error with formatted message
func NewAPIf(format string, args ...any) *Error {
	return New(WarnLevel, API, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// NETWORK ERRORS
// ============================================================================

// NewNetwork creates a new network error (WARN level)
func NewNetwork(message string, cause error) *Error {
	return New(WarnLevel, Network, message, cause)
}

// NewNetworkf creates a new network error with formatted message
func NewNetworkf(format string, args ...any) *Error {
	return New(WarnLevel, Network, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// AUTH ERRORS
// ============================================================================

// NewAuth creates a new authentication error (CRITICAL level)
func NewAuth(message string, cause error) *Error {
	return New(CriticalLevel, Auth, message, cause)
}

// NewAuthf creates a new authentication error with formatted message
func NewAuthf(format string, args ...any) *Error {
	return New(CriticalLevel, Auth, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// UNKNOWN ERRORS
// ============================================================================

// NewUnknown creates a new unknown error (CRITICAL level)
func NewUnknown(message string, cause error) *Error {
	return New(CriticalLevel, Unknown, message, cause)
}

// NewUnknownf creates a new unknown error with formatted message
func NewUnknownf(format string, args ...any) *Error {
	return New(CriticalLevel, Unknown, fmt.Sprintf(format, args...), nil)
}

// ============================================================================
// DEBUG ERRORS (for development/debugging purposes)
// ============================================================================

// NewDebug creates a new debug error (DEBUG level)
func NewDebug(message string, cause error) *Error {
	return New(DebugLevel, Unknown, message, cause)
}

// NewDebugf creates a new debug error with formatted message
func NewDebugf(format string, args ...any) *Error {
	return New(DebugLevel, Unknown, fmt.Sprintf(format, args...), nil)
}
