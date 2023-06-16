// Package util implements handy primitives that are too small to be placed
// inside their own separate package
package util

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

// FlagB is a convenience wrapper for boolean flag that picks its default value
// from the environment variable.
func FlagB(name string, defaultValue bool, usage string) *bool {
	value := defaultValue
	if raw := getEnv(name); raw != "" {
		value = IsTruthy(raw)
	}
	return flag.Bool(name, value, usage)
}

// FlagI is a convenience wrapper for integer flag that picks its default value
// from the environment variable.
func FlagI(name string, defaultValue int, usage string) *int {
	value := defaultValue
	if raw := getEnv(name); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			value = parsed
		}
	}
	return flag.Int(name, value, usage)
}

// FlagS is a convenience wrapper for string flag that picks its default value
// from the environment variable.
func FlagS(name, defaultValue, usage string) *string {
	value := defaultValue
	if raw := getEnv(name); raw != "" {
		value = raw
	}
	return flag.String(name, value, usage)
}

// getEnv converts name from kebab-case into UPPER_SNAKE_CASE and retrieves the
// environment variable with that name.
func getEnv(name string) string {
	env := strings.ToUpper(
		strings.ReplaceAll(name, "-", "_"),
	)
	return os.Getenv(env)
}
