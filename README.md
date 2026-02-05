📰 AI Tech News Agent

An AI-powered automation agent that collects, analyzes, summarizes, and delivers weekly technology news using Go, NewsAPI, Gemini AI, and Telegram Bot.

This project demonstrates a production-style backend agent that fetches news data, processes it with AI, and sends structured summaries automatically to Telegram.

🚀 Features
	•	📡 Fetches latest tech news from NewsAPI
	•	🤖 Summarizes articles using Google Gemini AI
	•	🧠 Extracts key topics & trending stories
	•	📰 Generates structured weekly tech summary
	•	📤 Sends formatted reports to Telegram
	•	⏱ Runs automatically via cron schedule
	•	🛡 Error handling & fallback support
	•	🧱 Clean, modular Go architecture

🛠 Tech Stack

Backend
	•	Go (Golang) → Core backend & agent logic
	•	REST APIs → External API integrations
	•	Cron Scheduler → Automated weekly execution

AI & Data
	•	Google Gemini API → News summarization & analysis
	•	NewsAPI → News data collection

Notifications
	•	Telegram Bot API → Message delivery
	•	Markdown formatted summary reports

  🧠 How It Works
	1.	Agent runs on scheduled cron job
	2.	Fetches news articles by category (AI, tech, etc.)
	3.	Filters & prepares article data
	4.	Sends content to Gemini AI for summarization
	5.	Generates:
	•	Executive summary
	•	Key topics
	•	Trending stories
	6.	Sends final formatted report to Telegram

tech-news-agent/
│
├── cmd/                # Application entrypoint
├── internal/
│   ├── collector/      # News fetching logic
│   ├── ai/             # Gemini summarization
│   ├── notifier/       # Telegram sender
│   ├── models/         # Data models
│   └── scheduler/      # Cron jobs
│
├── .env.example
├── go.mod
└── README.md

⚙️ Environment Variables

Create a .env file:

# TELEGRAM
TELEGRAM_CHAT_ID=<your_chat_id>
TELEGRAM_BOT_TOKEN=<your_bot_token>
# SCHEDULE
CRON_SCHEDULE=<your_schedule> --> e.g. 0 9 * * 1
#GEMINI
GEMINI_MODEL=<your_gemini_models>
GEMINI_API_KEY=<your_gemini_api_key>
#NEWS API
NEWS_API_KEY=<your_news_api_key>
MAX_NEWS_ARTICLES=<your_max_article>
