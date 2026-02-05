package services

import (
	"fmt"
	"strings"
	"tech-news-agent/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramNotifier handles sending notifications via Telegram
type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

// NewTelegramNotifier creates a new Telegram notifier instance
func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("creating Telegram bot: %w", err)
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

// SendSummary sends the news summary to the configured chat
func (tn *TelegramNotifier) SendSummary(summary *models.NewsSummary) error {
	message := tn.formatMessage(summary)

	// Split message if it's too long (Telegram limit is 4096 characters)
	messages := tn.splitMessage(message, 4000)

	for _, msg := range messages {
		telegramMsg := tgbotapi.NewMessage(tn.chatID, msg)
		telegramMsg.ParseMode = "Markdown"
		telegramMsg.DisableWebPagePreview = true

		if _, err := tn.bot.Send(telegramMsg); err != nil {
			return fmt.Errorf("sending message: %w", err)
		}
	}

	return nil
}

// SendError sends an error notification
func (tn *TelegramNotifier) SendError(errMsg string) error {
	message := fmt.Sprintf("⚠️ *Tech News Agent Error*\n\n```\n%s\n```", errMsg)
	msg := tgbotapi.NewMessage(tn.chatID, message)
	msg.ParseMode = "Markdown"

	if _, err := tn.bot.Send(msg); err != nil {
		return fmt.Errorf("sending error message: %w", err)
	}

	return nil
}

func (tn *TelegramNotifier) formatMessage(summary *models.NewsSummary) string {
	var sb strings.Builder

	// Header
	sb.WriteString("📰 *Weekly Tech News Summary*\n")
	sb.WriteString(fmt.Sprintf("📅 *%s*\n", summary.WeekRange))
	sb.WriteString(fmt.Sprintf("📊 Articles analyzed: %d\n", summary.TotalArticles))
	sb.WriteString("\n━━━━━━━━━━━━━━━━━\n\n")

	// Main Summary
	sb.WriteString(summary.Summary)
	sb.WriteString("\n\n")

	// Key Topics
	if len(summary.KeyTopics) > 0 {
		sb.WriteString("━━━━━━━━━━━━━━━━━\n")
		sb.WriteString("🔑 *Key Topics*\n\n")
		for _, topic := range summary.KeyTopics {
			sb.WriteString(fmt.Sprintf("• %s\n", topic))
		}
		sb.WriteString("\n")
	}

	// Trending Stories
	if len(summary.TrendingStories) > 0 {
		sb.WriteString("━━━━━━━━━━━━━━━━━\n")
		sb.WriteString("🔥 *Trending Stories*\n\n")
		for i, story := range summary.TrendingStories {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, story))
		}
		sb.WriteString("\n")
	}

	// Footer
	sb.WriteString("━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf("🤖 Generated on %s\n", summary.GeneratedAt.Format("Jan 02, 2006 15:04 MST")))
	sb.WriteString("_Powered by Gemini AI & Go_")

	return sb.String()
}

func (tn *TelegramNotifier) splitMessage(message string, maxLength int) []string {
	if len(message) <= maxLength {
		return []string{message}
	}

	var messages []string
	lines := strings.Split(message, "\n")
	var currentMsg strings.Builder

	for _, line := range lines {
		// If adding this line would exceed limit, start new message
		if currentMsg.Len()+len(line)+1 > maxLength {
			messages = append(messages, currentMsg.String())
			currentMsg.Reset()
		}

		if currentMsg.Len() > 0 {
			currentMsg.WriteString("\n")
		}
		currentMsg.WriteString(line)
	}

	if currentMsg.Len() > 0 {
		messages = append(messages, currentMsg.String())
	}

	return messages
}

// TestConnection sends a test message to verify the bot is working
func (tn *TelegramNotifier) TestConnection() error {
	msg := tgbotapi.NewMessage(tn.chatID, "✅ Tech News Agent is connected and ready!")
	if _, err := tn.bot.Send(msg); err != nil {
		return fmt.Errorf("test message failed: %w", err)
	}
	return nil
}
