package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tech-news-agent/internal/config"
	"tech-news-agent/internal/services"

	"github.com/robfig/cron/v3"
)

func main() {
	// Command line flags
	testMode := flag.Bool("test", false, "Run once immediately for testing")
	testConnection := flag.Bool("test-connection", false, "Test connections only")
	flag.Parse()

	// Initialize logger
	logger := log.New(os.Stdout, "[TechNewsAgent] ", log.LstdFlags|log.Lshortfile)

	logger.Println("🚀 Tech News Agent starting...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}
	logger.Printf("Configuration loaded successfully")
	logger.Printf("Using Gemini model: %s", cfg.GeminiModel)
	logger.Printf("Schedule: %s", cfg.CronSchedule)

	// Create news agent
	agent, err := services.NewNewsAgent(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to create news agent: %v", err)
	}
	defer agent.Close()

	// Test connection mode
	if *testConnection {
		logger.Println("Testing connections...")
		if err := agent.TestConnection(); err != nil {
			logger.Fatalf("Connection test failed: %v", err)
		}
		logger.Println("✅ Connection test passed!")
		return
	}

	// Test mode - run once immediately
	if *testMode {
		logger.Println("Running in test mode (single execution)...")
		if err := agent.TestRun(); err != nil {
			logger.Fatalf("Test run failed: %v", err)
		}
		logger.Println("Test run completed successfully!")
		return
	}

	//For check supported gemini models in SDK
	/*
		if *testMode {
			log.Println("🔍 Listing available Gemini models...")
			err := services.ListAvailableModels(cfg.GeminiAPIKey)
			if err != nil {
				log.Fatalf("Failed to list models: %v", err)
			}
			return
		}

	*/

	// Production mode - run on schedule
	logger.Printf("Running in production mode with schedule: %s", cfg.CronSchedule)

	// Create cron scheduler
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(logger)))

	_, err = c.AddFunc(cfg.CronSchedule, func() {
		logger.Println("Cron job triggered")
		ctx := context.Background()
		if err := agent.Run(ctx); err != nil {
			logger.Printf("❌ Job execution failed: %v", err)
		}
	})
	if err != nil {
		logger.Fatalf("Failed to add cron job: %v", err)
	}

	// Start the scheduler
	c.Start()
	logger.Println("✅ Scheduler started successfully")
	logger.Println("Press Ctrl+C to stop")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Println("Shutting down gracefully...")
	c.Stop()
	logger.Println("Goodbye!")
}
