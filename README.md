# telegram-go-chat-skeleton

Minimal Telegram bot skeleton in Go with long polling, configuration management, and optional Gemini API integration.

## Disclaimer

This is a personal skeleton project I use as a starting point for my own Telegram bots. It's not intended to be a universal framework, but you're welcome to use it, fork it, or steal ideas from it.

## Features

- **Long polling** via [telego](https://github.com/mymmrac/telego)
- **Configuration** via `.env` file and YAML prompts
- **Gemini API** integration (optional, can be removed)
- **Command routing** with a simple dispatcher
- **Graceful shutdown** via context cancellation

## Quick start

1. Clone the repo
2. Copy `.env.example` to `.env` and fill in your tokens
3. Copy `config.example.yaml` to `config.yaml`
4. Run:
```bash
go run cmd/bot/main.go
```

## Project structure

```
cmd/bot/main.go        # entry point
internal/
  bot/                 # bot setup, handlers, dispatcher
  config/              # .env and YAML config loader
  gemini/              # Gemini API client (optional)
prompts.example.yaml   # example prompts file
.env.example           # example environment file
```

## Configuration

- `TELEGRAM_BOT_TOKEN` — your bot token from BotFather
- `TELEGRAM_CHAT_ID` — target chat ID (group or private)
- `GEMINI_API_KEY` — Gemini API key (optional)

## Commands

- `/start` — welcome message
- `/news` — generate a news article via Gemini (if configured)

## License

MIT
