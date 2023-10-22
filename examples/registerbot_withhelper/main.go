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
