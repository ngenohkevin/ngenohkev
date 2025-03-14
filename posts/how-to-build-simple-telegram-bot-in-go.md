date:2025-03-08
# How to Build a Simple Telegram Bot in Golang
![Webserver Illustration](/static/img/go_bot.png)
This tutorial will guide you through building a simple telegram bot using Go.
## Prerequisites

Before starting, ensure you have the following:

- Go installed on your system ([Download Go](https://go.dev/dl/))
- A Telegram account
- A Telegram bot token (get one via [BotFather](https://t.me/BotFather))

## Step 1: Create a New Telegram Bot

1. Open Telegram and search for `@BotFather`.
2. Start a chat and type `/newbot`.
3. Follow the instructions to set up the bot.
4. Copy and save the bot token provided by BotFather.

## Step 2: Create a New Go Project

```sh
mkdir telegram-bot
cd telegram-bot
go mod init telegram-bot
```

## Step 3: Install Dependencies

We'll use the `github.com/go-telegram-bot-api/telegram-bot-api/v5` package to interact with the Telegram API.

```sh
go get github.com/go-telegram-bot-api/telegram-bot-api/v5
```

## Step 4: Write the Bot Code

Create a new file `main.go` and add the following code:

```go
package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN") // Set this in your environment variables
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // Ignore non-message updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! You said: "+update.Message.Text)
		bot.Send(msg)
	}
}
```

## Step 5: Set Up Environment Variable

Export your bot token as an environment variable:

```sh
export TELEGRAM_BOT_TOKEN="your-telegram-bot-token"
```

Or create a `.env` file and load it using a package like `github.com/joho/godotenv`.

## Step 6: Run the Bot

Execute the bot with:

```sh
go run main.go
```

## Step 7: Test Your Bot

- Open Telegram
- Search for your bot
- Start a chat and send a message
- The bot should reply with your message prefixed by `Hello! You said:`

## Conclusion

Congratulations! You've built a simple Telegram bot in Golang. You can extend it by handling commands, integrating APIs, and more.
