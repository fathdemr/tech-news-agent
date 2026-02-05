package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey     string
	TelegramBotToken string
	TelegramChatID   int64
	NewsAPIKey       string
	CronSchedule     string
	MaxNewsArticles  int
	GeminiModel      string
	NewsCategories   []string
}

// Load holds all application configuration
func Load() (*Config, error) {

	_ = godotenv.Load()

	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		return nil, err
	}

	maxArticles := 20
	if max := os.Getenv("MAX_NEWS_ARTICLES"); max != "" {
		if parsed, err := strconv.Atoi(max); err == nil {
			maxArticles = parsed
		}
	}

	cronSchedule := os.Getenv("CRON_SCHEDULE")
	if cronSchedule != "" {
		//Default: Every Monday at 9 AM
		cronSchedule = "0 9 * * 1"
	}

	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiModel != "" {
		geminiModel = "gemini-2.5-flash"
	}

	cfg := &Config{
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   chatID,
		NewsAPIKey:       os.Getenv("NEWS_API_KEY"),
		CronSchedule:     cronSchedule,
		MaxNewsArticles:  maxArticles,
		GeminiModel:      geminiModel,
		NewsCategories:   []string{"technology", "science", "business"},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required configurations is present
func (c *Config) Validate() error {
	if c.GeminiAPIKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is required")
	}
	if c.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if c.TelegramChatID == 0 {
		return fmt.Errorf("TELEGRAM_CHAT_ID is required")
	}
	if c.NewsAPIKey == "" {
		return fmt.Errorf("NEWS_API_KEY is required")
	}
	return nil
}
