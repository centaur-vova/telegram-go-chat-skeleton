// Package profanity provides a plugin that filters forbidden words.
//
// It checks incoming messages against a configurable list of bad words.
// If a match is found, the message is deleted and a warning is sent to the chat.
package profanity

import (
	"context"
	"strings"

	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/config"
	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/logger"
	"github.com/mymmrac/telego"
)

// Plugin implements the profanity filter.
type Plugin struct {
	cfg *config.ProfanityConfig
}

// New creates a new profanity filter plugin instance.
func New(cfg *config.ProfanityConfig) *Plugin {
	return &Plugin{cfg: cfg}
}

// Name returns the plugin identifier.
func (p *Plugin) Name() string {
	return "profanity"
}

// Handle checks the message for bad words.
//
// If a bad word is detected, it attempts to delete the message and
// sends a warning to the chat. Returns true if the message was handled
// (i.e., it contained profanity), preventing further processing.
func (p *Plugin) Handle(ctx context.Context, bot *telego.Bot, update *telego.Update) bool {
	if !p.cfg.Enabled {
		return false
	}

	if update.Message == nil {
		return false
	}

	text := strings.ToLower(update.Message.Text)
	for _, word := range p.cfg.BadWords {
		if strings.Contains(text, strings.ToLower(word)) {
			err := bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				MessageID: update.Message.MessageID,
			})
			if err != nil {
				logger.Error("Failed to delete message", "error", err)
			}

			_, err = bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "Космофлот не одобряет мат 🚀", // TODO change verbiage
			})
			if err != nil {
				logger.Error("Failed to send warning", "error", err)
			}

			return true
		}
	}

	return false
}
