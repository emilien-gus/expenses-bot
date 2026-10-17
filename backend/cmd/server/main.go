package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-bot-expenses/backend/internal/bot"
	"github.com/ai-bot-expenses/backend/internal/config"
	"github.com/ai-bot-expenses/backend/internal/db"
	"github.com/ai-bot-expenses/backend/internal/handler"
	"github.com/ai-bot-expenses/backend/internal/llm"
	"github.com/ai-bot-expenses/backend/internal/repository"
	"github.com/ai-bot-expenses/backend/internal/service"
)

func main() {
	log.Println("Starting Expense Tracking Service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	database, err := db.Init(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	expenseRepo := repository.NewExpenseRepository(database)
	expenseService := service.NewExpenseService(expenseRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	router := handler.SetupRouter(expenseHandler)

	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Root context cancelled on termination signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize LLM and Telegram Bot
	llmClient := llm.NewClient(cfg)
	tgBot, err := bot.NewBot(cfg, expenseService, llmClient)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram Bot: %v", err)
	}

	// Start Telegram bot listener in goroutine
	go tgBot.Start(ctx)

	// Start HTTP server in goroutine
	go func() {
		log.Printf("HTTP server listening on port %s", cfg.ServerPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")

	// Cancel bot context
	cancel()

	// Shutdown HTTP server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced shutdown: %v", err)
	}

	log.Println("Service exited properly.")
}
