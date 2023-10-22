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
