package logging

import (
	"os"

	logpkg "github.com/charmbracelet/log"
)

// Init initializes the package-level logger. Call once at program startup.
func Init() {
	// ensure default logger writes to stdout
	logpkg.SetOutput(os.Stdout)
}

// Info logs an informational message with optional key/value pairs.
func Info(msg string, kv ...any) {
	logpkg.Info(msg, kv...)
}

// Warn logs a warning message with optional key/value pairs.
func Warn(msg string, kv ...any) {
	logpkg.Warn(msg, kv...)
}

// Error logs an error message with optional key/value pairs.
func Error(msg string, kv ...any) {
	logpkg.Error(msg, kv...)
}

// Fatal logs a fatal error and exits.
func Fatal(msg string, kv ...any) {
	logpkg.Fatal(msg, kv...)
}

// Success logs a success message with a checkmark icon and optional key/value pairs.
// It uses Info level but marks the entry with a `status=success` key.
func Success(msg string, kv ...any) {
	// append status key so consumers can filter on it
	kv = append(kv, "status", "success")
	logpkg.Info("✓ "+msg, kv...)
}

// Step logs a progress/step message with an arrow icon and optional key/value pairs.
func Step(msg string, kv ...any) {
	logpkg.Info("→ "+msg, kv...)
}
