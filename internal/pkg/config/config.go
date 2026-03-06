package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultPort         = 50051
	defaultSpannerDB    = "projects/catalog-service/instances/catalog-instance/databases/catalog-db"
	defaultMigrationDir = "migrations"
)

type (
	Config struct {
		App      AppConfig      `mapstructure:"app"`
		Database DatabaseConfig `mapstructure:"database"`
	}

	AppConfig struct {
		Port         int    `mapstructure:"port"`
		MigrationDir string `mapstructure:"migrations_dir"`
	}

	DatabaseConfig struct {
		SpannerDB string `mapstructure:"spanner_db"`
	}
)

func InitConfig(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	setFromEnv(&cfg)
	return &cfg, nil
}

func setFromEnv(cfg *Config) {
	if port := os.Getenv("PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.App.Port)
	}
	if db := os.Getenv("SPANNER_DB"); db != "" {
		cfg.Database.SpannerDB = db
	}
	if mDir := os.Getenv("MIGRATIONS_DIR"); mDir != "" {
		cfg.App.MigrationDir = mDir
	}
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if env == "" || env == "local" {
		return nil
	}

	viper.SetConfigName(env)
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to merge env config file: %w", err)
		}
	}

	return nil
}

func populateDefaults() {
	viper.SetDefault("app.port", defaultPort)
	viper.SetDefault("app.migrations_dir", defaultMigrationDir)
	viper.SetDefault("database.spanner_db", defaultSpannerDB)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
