# TelePrompt
`TelePrompt` harnesses the capabilities of Channels and Goroutines to offer a seamless alternative to state machines and external storage when crafting forms and prompts for user data collection. Designed with customization in mind, this package accelerates your development process, enabling you to focus on building interactive and user-friendly Telegram bots with ease.
## Requirements
This package is a wrapper over `"github.com/go-telegram-bot-api/telegram-bot-api/v5"`

## Why TelePrompt?
- **Simplified Prompting**: Instead of manually managing each message, response, and error, TelePrompt encapsulates these into simple, reusable functions.
- **Error Handling**: Automatic error handling and retries without the need to manually set up each scenario.
- **Flexible & Extensible**: Built with generics in mind, TelePrompter can adapt to a variety of data types and can be easily extended for custom validation or prompts.
- **Open 2 Contribute**: With extension in mind, this package is open to new amazing methods.

## Usage
You can use this package either with its methods or with the helpers. You can read the examples in the `example` folder.
#### With Helpers
```go
package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tele_prompt "github.com/mhrlife/tele-prompt"
	"log"
	"os"
)

var TelegramToken = os.Getenv("TOKEN")

func main() {

	manager := tele_prompt.NewManager(tele_prompt.DefaultConfig())
	bot, err := tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if manager.Handle(update) {
			continue
		}
		go run(bot, manager, update)
	}

}

func run(bot *tgbotapi.BotAPI, manager *tele_prompt.Manager, update tgbotapi.Update) {
	userId := update.Message.From.ID

	name, err := tele_prompt.SendMessage[string](bot, tele_prompt.SendMessageConfig[string]{
		InitialMessage: tgbotapi.NewMessage(userId, "What is your name?"),
		Handler: func() (string, error) {
			return manager.TextField(context.Background(), userId)
		},
	})

	if err != nil {
		log.Println(err)
		return
	}

	age, err := tele_prompt.SendMessage[int](bot, tele_prompt.SendMessageConfig[int]{
		InitialMessage: tgbotapi.NewMessage(userId, fmt.Sprintf("Hi %s, What is your age?", name)),
		OnErrorMessage: tgbotapi.NewMessage(userId, "Your name must be more than 18.\n\nHow old are you?"),
		Handler: func() (int, error) {
			return manager.Range(context.Background(), userId, 18, 120)
		},
		MaxRetry: 3,
	})
	if err != nil {
		log.Println(err)
		return
	}
	bot.Send(tgbotapi.NewMessage(userId, fmt.Sprintf("%s, %d. Completed!", name, age)))
}
```
#### Without Helpers
```go
package main

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tele_prompt "github.com/mhrlife/tele-prompt"
	"log"
	"os"
)

var TelegramToken = os.Getenv("TOKEN")

func main() {

	manager := tele_prompt.NewManager(tele_prompt.DefaultConfig())
	bot, err := tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if manager.Handle(update) {
			continue
		}
		go run(bot, manager, update)
	}

}

func run(bot *tgbotapi.BotAPI, manager *tele_prompt.Manager, update tgbotapi.Update) {
	userId := update.Message.From.ID
	bot.Send(tgbotapi.NewMessage(userId, "What is your name?"))
	name, err := manager.TextField(context.Background(), userId)
	if err != nil {
		log.Println(err)
		return
	}
	bot.Send(tgbotapi.NewMessage(userId, fmt.Sprintf("Hi %s, What is your age?", name)))
	age, err := manager.Range(context.Background(), userId, 18, 120)
	if err != nil {
		// bad user response
		if errors.Is(err, tele_prompt.ErrOutofRange) {
			bot.Send(tgbotapi.NewMessage(userId, "you must be at least 18 to use this bot"))
			return
		}
		log.Println(err)
		return
	}
	bot.Send(tgbotapi.NewMessage(userId, fmt.Sprintf("%s, %d. Completed!", name, age)))
}

```
## Package Methods
- **Text, Image, Voice, Video, VideoNote, Document**: Waits for and retrieves specific types of media from the user.
- **Range**: Asks the user for a numerical input within a specified range.
- **OneOf**: Ensures that the user's response matches one of the provided options.
- **Url**: Validates if the user's input is a well-formed URL.
- **Email**: Checks if the user's input is a valid email address.

## Helpers
To use the SendMessage function, you need to provide a configuration using SendMessageConfig, which dictates the messaging behavior, error handling, and data retrieval.
```go
result, err := tele_prompt.SendMessage[Type](bot, tele_prompt.SendMessageConfig[Type]{
    InitialMessage:    tgbotapi.NewMessage(userId, "Your prompt message here"),
    OnErrorMessage:    tgbotapi.NewMessage(userId, "Message to send on error"), // optional
    OnMaxRetryMessage: tgbotapi.NewMessage(userId, "Message to send after maximum retries"), // optional
    Handler: func() (Type, error) {
        // Your handler logic here, you can use the manager.Method here
    },
    MaxRetry: 3, // Number of retries before giving up, optional
})
```


## Getting Started
### Installation
To get started with TelePrompt, simply add it to your Go project:
```bash
go get github.com/mhrlife/tele-prompt
```
### Setting up Updates Handling
In the section of your code where updates are managed, ensure that each update is first passed to TelePrompt's `Handle` method. If `Handle` processes the update, it returns true, indicating that the update should be skipped for further processing.
```go
updates := bot.GetUpdatesChan(u)

for update := range updates {
    if manager.Handle(update) {
        continue
    }
    go run(bot, manager, update)
}
```
By following the above pattern, TelePrompt seamlessly integrates into your bot's logic, efficiently managing prompts and user responses.
## Open to Contribution
We believe in the collective power of the open-source community. If you have ideas, bug fixes, or enhancements for TelePrompt, we welcome your contributions!
