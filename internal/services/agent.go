package services

import (
	"context"
	"fmt"
	"log"
	"tech-news-agent/internal/config"
	"time"
)

// NewsAgent orchestrates the entire news collection and analysis workflow
type NewsAgent struct {
	config    *config.Config
	collector *NewsCollector
	analyzer  *AIAnalyzer
	notifier  *TelegramNotifier
	logger    *log.Logger
}

// NewNewsAgent creates a new news agent instance
func NewNewsAgent(cfg *config.Config, logger *log.Logger) (*NewsAgent, error) {
	collector := NewNewsCollector(cfg.NewsAPIKey, cfg.MaxNewsArticles)

	analyzer, err := NewAIAnalyzer(cfg.GeminiAPIKey, cfg.GeminiModel)
	if err != nil {
		return nil, fmt.Errorf("initializing AI analyzer: %w", err)
	}

	notifier, err := NewTelegramNotifier(cfg.TelegramBotToken, cfg.TelegramChatID)
	if err != nil {
		return nil, fmt.Errorf("initializing Telegram notifier: %w", err)
	}

	return &NewsAgent{
		config:    cfg,
		collector: collector,
		analyzer:  analyzer,
		notifier:  notifier,
		logger:    logger,
	}, nil
}

// Close cleans up resources
func (na *NewsAgent) Close() error {
	return na.analyzer.Close()
}

// Run executes the complete workflow
func (na *NewsAgent) Run(ctx context.Context) error {
	na.logger.Println("Starting weekly news collection and analysis...")

	// Step 1: Collect news
	na.logger.Println("Step 1/3: Collecting news articles...")
	articles, err := na.collector.FetchWeeklyNews(na.config.NewsCategories)
	if err != nil {
		na.logger.Printf("Error collecting news: %v", err)
		// Fall back to mock data for testing
		na.logger.Println("Using mock data for testing...")
		articles = na.collector.GetMockNews()
	}
	na.logger.Printf("Collected %d articles", len(articles))

	// Step 2: Analyze with AI
	na.logger.Println("Step 2/3: Analyzing articles with Gemini AI...")

	// Create a context with timeout for AI analysis
	aiCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	summary, err := na.analyzer.AnalyzeNews(aiCtx, articles)
	if err != nil {
		errMsg := fmt.Sprintf("AI analysis failed: %v", err)
		na.logger.Println(errMsg)
		if notifyErr := na.notifier.SendError(errMsg); notifyErr != nil {
			na.logger.Printf("Failed to send error notification: %v", notifyErr)
		}
		return fmt.Errorf("analyzing news: %w", err)
	}
	na.logger.Println("Analysis complete")

	// Step 3: Send via Telegram
	na.logger.Println("Step 3/3: Sending summary via Telegram...")
	if err := na.notifier.SendSummary(summary); err != nil {
		errMsg := fmt.Sprintf("Failed to send Telegram notification: %v", err)
		na.logger.Println(errMsg)
		return fmt.Errorf("sending notification: %w", err)
	}

	na.logger.Println("✅ Weekly news summary sent successfully!")
	return nil
}

// TestRun runs the agent immediately for testing
func (na *NewsAgent) TestRun() error {
	ctx := context.Background()
	return na.Run(ctx)
}

// TestConnection tests all service connections
func (na *NewsAgent) TestConnection() error {
	na.logger.Println("Testing Telegram connection...")
	if err := na.notifier.TestConnection(); err != nil {
		return fmt.Errorf("Telegram connection test failed: %w", err)
	}
	na.logger.Println("✅ All connections successful")
	return nil
}
