package models

import "time"

type Article struct {
	Title       string    `json:"title"`
	Desc        string    `json:"desc"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"publishedAt"`
	Category    string    `json:"category"`
}

type NewsSummary struct {
	WeekRange       string    `json:"weekRange"`
	TotalArticles   int       `json:"totalArticles"`
	Summary         string    `json:"summary"`
	KeyTopics       []string  `json:"keyTopics"`
	TrendingStories []string  `json:"trendingStories"`
	GeneratedAt     time.Time `json:"generatedAt"`
}
