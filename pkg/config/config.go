package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DatabaseConfig holds configuration for the database connection.
type DatabaseConfig struct {
	Driver          string        `yaml:"driver" json:"driver"`
	Host            string        `yaml:"host" json:"host"`
	Port            int           `yaml:"port" json:"port"`
	Database        string        `yaml:"database" json:"database"`
	Username        string        `yaml:"username" json:"username"`
	Password        string        `yaml:"password" json:"password"`
	SSLMode         string        `yaml:"ssl_mode" json:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" json:"conn_max_idle_time"`
	BusyTimeout     int           `yaml:"busy_timeout" json:"busy_timeout"`
	ConnectTimeout  int           `yaml:"connect_timeout" json:"connect_timeout"`
	AutoMigrate     bool          `yaml:"auto_migrate" json:"auto_migrate"`
}

// LoggerConfig holds configuration for the logger.
type LoggerConfig struct {
	Level         string `yaml:"level" json:"level"`
	Format        string `yaml:"format" json:"format"`
	EnableConsole bool   `yaml:"enable_console" json:"enable_console"`
	ConsoleToStd  string `yaml:"console_to_std" json:"console_to_stdout"`
	EnableFile    bool   `yaml:"enable_file" json:"enable_file"`
	FilePath      string `yaml:"file_path" json:"file_path"`
	FileMaxAge    int    `yaml:"file_max_age" json:"file_max_age"`
	Development   bool   `yaml:"development" json:"development"`
}

// AppSectionConfig holds general application configuration.
type AppSectionConfig struct {
	Name        string `yaml:"name" json:"name"`
	Version     string `yaml:"version" json:"version"`
	Environment string `yaml:"environment" json:"environment"`
}

// AuditConfig holds audit configuration.
type AuditConfig struct {
	AlertWebhook string `yaml:"alert_webhook" json:"alert_webhook"`
	ArchiveDir   string `yaml:"archive_dir" json:"archive_dir"`
	Retention    string `yaml:"retention" json:"retention"`
}

// BrainConfig holds brain configuration.
type BrainConfig struct {
	Enabled             bool `yaml:"enabled" json:"enabled"`
	DisableOptimization bool `yaml:"disable_optimization" json:"disable_optimization"`
}

// ServerConfig holds server configuration.
type ServerConfig struct {
	JWTSecret string        `yaml:"jwt_secret" json:"jwt_secret"`
	JWTExpiry time.Duration `yaml:"jwt_expiry" json:"jwt_expiry"`
	Port      int           `yaml:"port" json:"port"`
}

// AppConfig holds the global configuration.
type AppConfig struct {
	App      AppSectionConfig `yaml:"app"`
	Server   ServerConfig     `yaml:"server"`
	Database DatabaseConfig   `yaml:"database"`
	Logger   LoggerConfig     `yaml:"logger"`
	Audit    AuditConfig      `yaml:"audit"`
	Brain    BrainConfig      `yaml:"brain"`
}

// LoadConfig loads configuration from a YAML file.
// Supports environment variable override: DATABASE_HOST → database.host
func LoadConfig(path string) (*AppConfig, error) {
	if path == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Environment variable override (e.g. DATABASE_HOST overrides config.database.host)
	overrideFromEnv(&config)

	return &config, nil
}

// Load is an alias for LoadConfig.
func Load(path string) (*AppConfig, error) {
	return LoadConfig(path)
}

func overrideFromEnv(cfg *AppConfig) {
	if v := os.Getenv("DATABASE_DRIVER"); v != "" {
		cfg.Database.Driver = v
	}
	if v := os.Getenv("DATABASE_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DATABASE_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Database.Port)
	}
	if v := os.Getenv("DATABASE_DATABASE"); v != "" {
		cfg.Database.Database = v
	}
	if v := os.Getenv("DATABASE_USERNAME"); v != "" {
		cfg.Database.Username = v
	}
	if v := os.Getenv("DATABASE_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("LOGGER_LEVEL"); v != "" {
		cfg.Logger.Level = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.Port)
	}
	if v := os.Getenv("SERVER_JWT_SECRET"); v != "" {
		cfg.Server.JWTSecret = v
	}
	// Convert "." notation env vars (from viper-style) to underscore
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(kv[0])
		key = strings.ReplaceAll(key, ".", "_")
		_ = key
	}
}
