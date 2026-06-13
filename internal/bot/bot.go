package bot

import (
	"context"
	"log"

	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/config"
	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/gemini"
	"github.com/mymmrac/telego"
)

type Bot struct {
	bot    *telego.Bot
	gemini *gemini.Client
	cfg    config.Config
}

func New(cfg config.Config) *Bot {
	// Create bot instance
	bot, err := telego.NewBot(cfg.TelegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	// Gemini client
	client := gemini.New(cfg.GeminiModel, cfg.GeminiKey)

	return &Bot{
		bot:    bot,
		gemini: client,
		cfg:    cfg,
	}
}

func (b *Bot) Subscribe(ctx context.Context) (<-chan telego.Update, error) {
	// Get updates channel via long polling with context
	params := &telego.GetUpdatesParams{
		Timeout: b.cfg.PollTimeout,
	}
	return b.bot.UpdatesViaLongPolling(ctx, params)
}

func (b *Bot) Handle(ctx context.Context, updates <-chan telego.Update) {
	chatID := b.cfg.ChatID
	prompts := b.cfg.Prompts

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Ignore messages from other chats
		if update.Message.Chat.ID != chatID.ID {
			continue
		}

		text := update.Message.Text

		// /news command — generate a real news article, well, "almost" real haha
		if text == "/news" {
			b.news(ctx, chatID, prompts)
		}

		// /start command — welcome message
		if text == "/start" {
			b.start(ctx, chatID, prompts)
		}
	}
}
