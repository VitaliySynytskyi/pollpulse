package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnv gets an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvInt gets an environment variable as int or returns a default value
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetEnvBool gets an environment variable as bool or returns a default value
func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// GetEnvDuration gets an environment variable as time.Duration or returns a default value
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	durationValue, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return durationValue
}

// GetEnvSlice gets an environment variable as a slice of strings or returns a default value
func GetEnvSlice(key, separator string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return strings.Split(value, separator)
}

// MustGetEnv gets an environment variable or panics if it's not set
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is required but not set", key))
	}
	return value
}
