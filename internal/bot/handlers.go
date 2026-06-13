package bot

import (
	"context"

	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/config"
	"github.com/centaur-vova/telegram-go-chat-skeleton/internal/logger"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func (b *Bot) news(ctx context.Context, chatID telego.ChatID, prompts *config.PromptsConfig) {
	msg, _ := b.bot.SendMessage(ctx, telegoutil.Message(chatID, prompts.Messages.Recording))

	userPrompt := `Вспомни одно реальное позитивное событие из жизни любой российской сельской глубинки за последние два года (например: благоустройство сквера, ремонт дороги, открытие почты, победа на районных соревнованиях, закупка тракторов). Перенеси это событие в р.п. Сурское и перепиши строго по инструкции`
	news, err := b.gemini.Ask(prompts.News.System, userPrompt, false)
	if err != nil {
		logger.Error("Gemini error", "error", err)
		b.bot.SendMessage(ctx, telegoutil.Message(chatID, prompts.Messages.Error))
		return
	}

	// Delete "recording..." message and send the news
	b.bot.DeleteMessage(ctx, telegoutil.Delete(chatID, msg.MessageID))
	b.bot.SendMessage(ctx, telegoutil.Message(chatID, "📡 "+news))
}

func (b *Bot) start(ctx context.Context, chatID telego.ChatID, prompts *config.PromptsConfig) {
	b.bot.SendMessage(ctx, telegoutil.Message(chatID, prompts.Messages.Welcome))
}
