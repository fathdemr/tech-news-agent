package services

import (
	"context"
	"fmt"
	"strings"
	"tech-news-agent/internal/models"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AIAnalyzer handles AI-powered news analysis using Gemini
type AIAnalyzer struct {
	client *genai.Client
	model  string
}

// NewAIAnalyzer creates a new AI analyzer instance
func NewAIAnalyzer(apiKey, modelName string) (*AIAnalyzer, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	return &AIAnalyzer{
		client: client,
		model:  modelName,
	}, nil
}

// Close closes the AI client connection
func (a *AIAnalyzer) Close() error {
	return a.client.Close()
}

// AnalyzeNews generates a comprehensive summary of news articles
func (a *AIAnalyzer) AnalyzeNews(ctx context.Context, articles []models.Article) (*models.NewsSummary, error) {
	if len(articles) == 0 {
		return nil, fmt.Errorf("no articles to analyze")
	}

	prompt := a.buildPrompt(articles)

	model := a.client.GenerativeModel(a.model)

	// Configure the model
	model.SetTemperature(0.7)
	model.SetTopP(0.9)
	model.SetTopK(40)
	model.SetMaxOutputTokens(2048)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("generating content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	summary := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// Extract key topics and trending stories from the summary
	keyTopics, trendingStories := a.extractInsights(summary)

	weekRange := fmt.Sprintf("%s - %s",
		time.Now().AddDate(0, 0, -7).Format("Jan 02"),
		time.Now().Format("Jan 02, 2006"))

	return &models.NewsSummary{
		WeekRange:       weekRange,
		TotalArticles:   len(articles),
		Summary:         summary,
		KeyTopics:       keyTopics,
		TrendingStories: trendingStories,
		GeneratedAt:     time.Now(),
	}, nil
}

func (a *AIAnalyzer) buildPrompt(articles []models.Article) string {
	var sb strings.Builder

	sb.WriteString("You are a professional tech news analyst. Analyze the following technology news articles from the past week and create a comprehensive weekly summary.\n\n")
	sb.WriteString("Articles:\n\n")

	for i, article := range articles {
		sb.WriteString(fmt.Sprintf("%d. Title: %s\n", i+1, article.Title))
		sb.WriteString(fmt.Sprintf("   Source: %s\n", article.Source))
		sb.WriteString(fmt.Sprintf("   Category: %s\n", article.Category))
		if article.Desc != "" {
			sb.WriteString(fmt.Sprintf("   Description: %s\n", article.Desc))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nPlease provide:\n")
	sb.WriteString("1. A concise executive summary (2-3 paragraphs) of the week's most important tech developments\n")
	sb.WriteString("2. Key topics and themes (list 3-5 main topics)\n")
	sb.WriteString("3. Top 3 trending stories with brief explanations\n")
	sb.WriteString("4. Notable insights or patterns across the news\n\n")
	sb.WriteString("Format your response in a clear, professional manner suitable for a weekly newsletter.\n")
	sb.WriteString("Use markdown formatting with headers (##) for sections.\n")

	return sb.String()
}

func (a *AIAnalyzer) extractInsights(summary string) ([]string, []string) {
	// Simple extraction logic - in production, you might use more sophisticated parsing
	keyTopics := []string{}
	trendingStories := []string{}

	lines := strings.Split(summary, "\n")
	inKeyTopics := false
	inTrendingStories := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(strings.ToLower(line), "key topics") ||
			strings.Contains(strings.ToLower(line), "main topics") {
			inKeyTopics = true
			inTrendingStories = false
			continue
		}

		if strings.Contains(strings.ToLower(line), "trending") ||
			strings.Contains(strings.ToLower(line), "top") {
			inTrendingStories = true
			inKeyTopics = false
			continue
		}

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if inKeyTopics && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•")) {
			topic := strings.TrimPrefix(line, "-")
			topic = strings.TrimPrefix(topic, "*")
			topic = strings.TrimPrefix(topic, "•")
			topic = strings.TrimSpace(topic)
			if topic != "" && len(keyTopics) < 5 {
				keyTopics = append(keyTopics, topic)
			}
		}

		if inTrendingStories && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•") ||
			(len(line) > 0 && line[0] >= '0' && line[0] <= '9')) {
			story := strings.TrimPrefix(line, "-")
			story = strings.TrimPrefix(story, "*")
			story = strings.TrimPrefix(story, "•")
			// Remove leading numbers and dots
			if len(story) > 0 && story[0] >= '0' && story[0] <= '9' {
				parts := strings.SplitN(story, ".", 2)
				if len(parts) > 1 {
					story = parts[1]
				}
			}
			story = strings.TrimSpace(story)
			if story != "" && len(trendingStories) < 3 {
				trendingStories = append(trendingStories, story)
			}
		}
	}

	// Fallback to generic topics if extraction failed
	if len(keyTopics) == 0 {
		keyTopics = []string{"Artificial Intelligence", "Cloud Computing", "Cybersecurity"}
	}
	if len(trendingStories) == 0 {
		trendingStories = []string{"Major tech industry developments", "Innovation breakthroughs", "Market trends"}
	}

	return keyTopics, trendingStories
}

func ListAvailableModels(apiKey string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	defer client.Close()

	it := client.ListModels(ctx)

	fmt.Println("📦 Available models:")
	for {
		model, err := it.Next()
		if err != nil {
			break
		}

		fmt.Printf("- %s\n", model.Name)
		fmt.Printf("  Supported methods: %v\n\n", model.SupportedGenerationMethods)
	}

	return nil
}
