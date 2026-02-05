package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tech-news-agent/internal/models"
	"time"
)

// NewsCollector handles fetching news from various sources
type NewsCollector struct {
	apiKey     string
	maxResults int
	httpClient *http.Client
}

// NewsAPIResponse represents the response from NewsAPI
type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			Name string `json:"name"`
		} `json:"source"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PublishedAt string `json:"publishedAt"`
	} `json:"articles"`
}

// NewNewsCollector creates a new news collector instance
func NewNewsCollector(apiKey string, maxResults int) *NewsCollector {
	return &NewsCollector{
		apiKey:     apiKey,
		maxResults: maxResults,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchWeeklyNews retrieves technology news from the past week
func (nc *NewsCollector) FetchWeeklyNews(categories []string) ([]models.Article, error) {
	var allArticles []models.Article

	for _, category := range categories {
		articles, err := nc.fetchByCategory(category)
		if err != nil {
			// Log error but continue with other categories
			fmt.Printf("Error fetching %s news: %v\n", category, err)
			continue
		}
		allArticles = append(allArticles, articles...)
	}

	if len(allArticles) == 0 {
		return nil, fmt.Errorf("no articles found")
	}

	return allArticles, nil
}

func (nc *NewsCollector) fetchByCategory(category string) ([]models.Article, error) {
	// Calculate date range (last 7 days)
	to := time.Now()
	from := to.AddDate(0, 0, -7)

	// Build NewsAPI URL
	baseURL := "https://newsapi.org/v2/everything"
	params := url.Values{}
	params.Add("q", category)
	params.Add("from", from.Format("2006-01-02"))
	params.Add("to", to.Format("2006-01-02"))
	params.Add("sortBy", "popularity")
	params.Add("language", "en")
	params.Add("pageSize", fmt.Sprintf("%d", nc.maxResults/len([]string{"politic", "sport", "business"})))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", nc.apiKey)

	resp, err := nc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	articles := make([]models.Article, 0, len(apiResp.Articles))
	for _, a := range apiResp.Articles {
		publishedAt, _ := time.Parse(time.RFC3339, a.PublishedAt)

		articles = append(articles, models.Article{
			Title:       a.Title,
			Desc:        a.Description,
			URL:         a.URL,
			Source:      a.Source.Name,
			PublishedAt: publishedAt,
			Category:    category,
		})
	}

	return articles, nil
}

// GetMockNews returns mock news for testing without API key
func (nc *NewsCollector) GetMockNews() []models.Article {
	return []models.Article{
		{
			Title:       "AI Breakthrough: New Language Model Surpasses Human Performance",
			Desc:        "Researchers announce a groundbreaking AI model that demonstrates superior performance across multiple benchmarks.",
			URL:         "https://example.com/ai-breakthrough",
			Source:      "TechCrunch",
			PublishedAt: time.Now().AddDate(0, 0, -1),
			Category:    "technology",
		},
		{
			Title:       "Quantum Computing Reaches New Milestone",
			Desc:        "Scientists achieve quantum supremacy with a 1000-qubit processor.",
			URL:         "https://example.com/quantum",
			Source:      "MIT Technology Review",
			PublishedAt: time.Now().AddDate(0, 0, -2),
			Category:    "science",
		},
		{
			Title:       "Major Tech Companies Announce Climate Initiatives",
			Desc:        "Leading technology firms commit to carbon neutrality by 2030.",
			URL:         "https://example.com/climate",
			Source:      "Bloomberg",
			PublishedAt: time.Now().AddDate(0, 0, -3),
			Category:    "business",
		},
	}
}
