package config

import (
	"os"
	"strconv"
)

type BotConfig struct {
	TelegramToken string
	Url           string
	ReportChatId  int64
}

type Config struct {
	BotConfig BotConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		BotConfig: BotConfig{
			TelegramToken: getEnv("TELEGRAM_TOKEN", ""),
			Url:           getEnv("URL", ""),
			ReportChatId:  getEnvAsInt("REPORT_CHAT_ID", 0),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}

	return defaultVal
}
