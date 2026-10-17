package bot

import (
	"context"
	"log"

	"github.com/ai-bot-expenses/backend/internal/config"
	"github.com/ai-bot-expenses/backend/internal/llm"
	"github.com/ai-bot-expenses/backend/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the Telegram bot service.
type Bot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
}

// NewBot initializes a new Telegram bot instance.
func NewBot(cfg *config.Config, svc *service.ExpenseService, llmClient *llm.Client) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account @%s", api.Self.UserName)

	handler := NewHandler(svc, llmClient)

	return &Bot{
		api:     api,
		handler: handler,
	}, nil
}

// Start begins long-polling for updates and processing incoming messages until context is cancelled.
func (b *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	log.Println("Telegram bot started listening for updates...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Telegram bot stopping...")
			b.api.StopReceivingUpdates()
			return
		case update, ok := <-updates:
			if !ok {
				return
			}

			if update.Message == nil || update.Message.Text == "" {
				continue
			}

			go b.processMessage(ctx, update.Message)
		}
	}
}

func (b *Bot) processMessage(ctx context.Context, msg *tgbotapi.Message) {
	userID := msg.From.ID
	chatID := msg.Chat.ID
	text := msg.Text

	log.Printf("Received message from user %d: %q", userID, text)

	replyText := b.handler.HandleMessage(ctx, userID, text)

	reply := tgbotapi.NewMessage(chatID, replyText)
	if _, err := b.api.Send(reply); err != nil {
		log.Printf("Failed to send reply to chat %d: %v", chatID, err)
	}
}
