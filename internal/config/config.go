// Package config provides configuration loading from .env and YAML files.
//
// It handles environment variables for Telegram bot token, Gemini API keys,
// chat ID, poll timeout, and log level. Prompts for commands are loaded
// from a separate YAML file (prompts.yaml).
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

type SecretsConfig struct {
	GeminiKey     string
	TelegramToken string
}

type TelegramConfig struct {
	ChatID int64 `yaml:"chat_id"`
}

type GeminiConfig struct {
	Model string `yaml:"model"`
}

type BotConfig struct {
	LogLevel    string `yaml:"log_level"`
	PollTimeout int    `yaml:"poll_timeout"`
}

// Config holds all configuration for the bot.
type Config struct {
	Secrets  SecretsConfig  `yaml:"-"`
	Bot      BotConfig      `yaml:"bot"`
	Telegram TelegramConfig `yaml:"telegram"`
	Gemini   GeminiConfig   `yaml:"gemini"`
	Prompts  *PromptsConfig `yaml:"prompts"`
}

// Load reads .env file, parses environment variables, and loads prompts from YAML.
// It exits with fatal log if required variables are missing or invalid.
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatalf("GEMINI_API_KEY is not set")
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	// Load config
	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg.Secrets = SecretsConfig{
		GeminiKey:     geminiKey,
		TelegramToken: telegramToken,
	}

	return cfg
}

// PromptsConfig contains all prompt templates loaded from YAML.
type PromptsConfig struct {
	Messages  promptsConfigMessages   `yaml:"messages"`
	News      struct{ System string } `yaml:"news"`
	Interview struct{ System string } `yaml:"interview"`
}

// PromptsConfigMessages contains message templates for bot responses.
type promptsConfigMessages struct {
	Recording string `yaml:"recording"`
	Error     string `yaml:"error"`
	Welcome   string `yaml:"welcome"`
}

// loadConfig reads and parses the config YAML file at the given path.
func loadConfig(path string) (*Config, error) {
	// #nosec G304 — prompts file path is configured by the developer
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
