package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort              = "8080"
	defaultReadTimeout       = 5 * time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultIdleTimeout       = 60 * time.Second
	defaultShutdownTimeout   = 10 * time.Second
	defaultTokenTTL          = 24 * time.Hour
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	HTTP        HTTPConfig
	Auth        AuthConfig
}

type HTTPConfig struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

type AuthConfig struct {
	TokenTTL time.Duration
}

func Load() (Config, error) {
	var errs []error

	databaseURL, err := requiredEnv("DATABASE_URL")
	if err != nil {
		errs = append(errs, err)
	}

	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		errs = append(errs, err)
	}

	port, err := stringEnv("PORT", defaultPort)
	if err != nil {
		errs = append(errs, err)
	}

	readTimeout, err := durationEnv("HTTP_READ_TIMEOUT", defaultReadTimeout)
	if err != nil {
		errs = append(errs, err)
	}

	readHeaderTimeout, err := durationEnv("HTTP_READ_HEADER_TIMEOUT", defaultReadHeaderTimeout)
	if err != nil {
		errs = append(errs, err)
	}

	writeTimeout, err := durationEnv("HTTP_WRITE_TIMEOUT", defaultWriteTimeout)
	if err != nil {
		errs = append(errs, err)
	}

	idleTimeout, err := durationEnv("HTTP_IDLE_TIMEOUT", defaultIdleTimeout)
	if err != nil {
		errs = append(errs, err)
	}

	shutdownTimeout, err := durationEnv("HTTP_SHUTDOWN_TIMEOUT", defaultShutdownTimeout)
	if err != nil {
		errs = append(errs, err)
	}

	tokenTTL, err := durationEnv("AUTH_TOKEN_TTL", defaultTokenTTL)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return Config{}, errors.Join(errs...)
	}

	return Config{
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		Port:        port,
		HTTP: HTTPConfig{
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
			ShutdownTimeout:   shutdownTimeout,
		},
		Auth: AuthConfig{
			TokenTTL: tokenTTL,
		},
	}, nil
}

func (c Config) Address() string {
	return ":" + c.Port
}

func requiredEnv(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}

func stringEnv(key, fallback string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		value = fallback
	}

	if _, err := strconv.Atoi(value); err != nil {
		return "", fmt.Errorf("%s must be a valid port: %w", key, err)
	}

	return value, nil
}

func durationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", key, err)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return duration, nil
}
