// Package bot implements the Telegram bot core functionality.
//
// It provides bot setup, long polling subscription, and command dispatching
// for /start and /news handlers.
package bot

import (
	"context"
	"log"

	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/config"
	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/gemini"
	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/plugins"
	"github.com/mymmrac/telego"
)

// Bot represents a Telegram bot instance with its dependencies.
type Bot struct {
	bot     *telego.Bot
	gemini  *gemini.Client
	cfg     *config.Config
	chatID  telego.ChatID
	plugins []plugins.Plugin
}

// New creates a new Bot instance. It initializes the Telegram bot and
// Gemini client. If bot creation fails, it logs a fatal error.
func New(cfg *config.Config, plugins []plugins.Plugin) *Bot {
	// Create bot instance
	bot, err := telego.NewBot(cfg.Secrets.TelegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	// Gemini client
	client := gemini.New(cfg.Gemini.Model, cfg.Secrets.GeminiKey)

	return &Bot{
		bot:     bot,
		gemini:  client,
		chatID:  telego.ChatID{ID: cfg.Telegram.ChatID},
		cfg:     cfg,
		plugins: plugins,
	}
}

// Subscribe starts long polling for Telegram updates and returns a channel
// of updates. The channel is closed when the context is cancelled.
func (b *Bot) Subscribe(ctx context.Context) (<-chan telego.Update, error) {
	// Get updates channel via long polling with context
	params := &telego.GetUpdatesParams{
		Timeout: b.cfg.Bot.PollTimeout,
	}
	return b.bot.UpdatesViaLongPolling(ctx, params)
}

// Handle processes incoming updates and dispatches commands to handlers.
// It filters messages from chats other than the configured chat ID.
func (b *Bot) Handle(ctx context.Context, updates <-chan telego.Update) {
	prompts := b.cfg.Prompts

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Ignore messages from other chats
		if update.Message.Chat.ID != b.chatID.ID {
			continue
		}

		// plugins
		handled := false
		for _, p := range b.plugins {
			if p.Handle(ctx, b.bot, &update) {
				handled = true
				break
			}
		}
		if handled {
			continue
		}

		text := update.Message.Text
		// /news command — generate a real news article, well, "almost" real haha
		if text == "/news" {
			b.news(ctx, prompts)
		}
		// /start command — welcome message
		if text == "/start" {
			b.start(ctx, prompts)
		}
	}
}
