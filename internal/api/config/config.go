// Package config provides application configuration management.
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"open-replays/internal/pkg/parse"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Video    VideoConfig    `mapstructure:"video"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Path        string        `mapstructure:"path"`
	BusyTimeout time.Duration `mapstructure:"busy-timeout"`
	JournalMode string        `mapstructure:"journal-mode"`
}

// StorageConfig holds file storage configuration.
type StorageConfig struct {
	BaseDir   string `mapstructure:"base-dir"`
	PublicURL string `mapstructure:"public-url"`
}

// VideoConfig holds video processing configuration.
type VideoConfig struct {
	WorkerCount       int      `mapstructure:"worker-count"`
	MaxFileSize       string   `mapstructure:"max-file-size"`
	AllowedExtensions []string `mapstructure:"allowed-extensions"`
}

// MaxFileSizeBytes returns the max file size in bytes.
func (v VideoConfig) MaxFileSizeBytes() (int64, error) {
	return parse.FileSizeToBytes(v.MaxFileSize)
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load loads configuration from environment variables and config files.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// Set defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.path", "db.sqlite3")
	viper.SetDefault("database.busy-timeout", "5s")
	viper.SetDefault("database.journal-mode", "WAL")
	viper.SetDefault("storage.base-dir", "uploads")
	viper.SetDefault("storage.public-url", "http://localhost:8080/media")
	viper.SetDefault("video.worker-count", 2)
	viper.SetDefault("video.allowed-extensions", []string{".mp4", ".webm", ".mov"})
	viper.SetDefault("video.max-file-size", "100MB")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// Environment variables
	viper.SetEnvPrefix("OPEN_REPLAYS")
	viper.AutomaticEnv()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError *viper.ConfigFileNotFoundError

		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("read config: %w", err)
		}
		// Config file not found is OK, we have defaults and env vars
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Address returns the server address in format "host:port".
func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// DSN returns the SQLite connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("file:%s?_busy_timeout=%d&_journal_mode=%s",
		d.Path,
		d.BusyTimeout.Milliseconds(),
		d.JournalMode)
}
