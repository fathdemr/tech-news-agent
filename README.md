# 📰 AI Tech News Agent

An AI-powered automation agent that collects, analyzes, summarizes, and delivers weekly technology news using Go, NewsAPI, Gemini AI, and Telegram Bot.

This project demonstrates a production-style backend agent that fetches news data, processes it with AI, and sends structured summaries automatically to Telegram.

---

## 🚀 Features
-	📡 Fetches latest tech news from NewsAPI
-	🤖 Summarizes articles using Google Gemini AI
-	🧠 Extracts key topics & trending stories
-	📰 Generates structured weekly tech summary
-	📤 Sends formatted reports to Telegram
-	⏱ Runs automatically via cron schedule
-	🛡 Error handling & fallback support
-	🧱 Clean, modular Go architecture

---

## 🛠️ Tech Stack

- **Go (Golang)** – BCore backend & agent logic
- **Google Gemini API** – News summarization & analysis
- **NewsAPI** –  News data collection
- **Telegram Bot API** – Message delivery
- **Postman** – API testing

---

## 🧠 How It Works
1.	Agent runs on scheduled cron job
2.	Fetches news articles by category (AI, tech, etc.)
3.	Filters & prepares article data
4.	Sends content to Gemini AI for summarization
5.	Generates:
•	Executive summary
•	Key topics
•	Trending stories
6.	Sends final formatted report to Telegram

---

### Clone Repository

```bash
git clone https://github.com/yourusername/tech-news-agent.git
cd tech-news-agent
```

Install dependencies:

```bash
go mod tidy
```

---

## ⚙️ Environment Variables

Create a .env file:
```env
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
```

---

## Running the Application

You can run the agent in **three different modes**.

---

## 🧪 Test Mode (run once instantly)

If you don’t want cron scheduling
and just want to run immediately:

```bash
go run cmd/main.go -test
```

---

## 🤖 Test Telegram Connection Only

To test Telegram bot connection:

```bash
go run cmd/main.go -test-connection
```


This sends a simple test message:

Tech News Agent is connected

---

## 🏗 Build Binary

```bash
go build -o news-agent cmd/main.go
```


Run:

```bash
./news-agent -test
```


or

./news-agent -testconnection
