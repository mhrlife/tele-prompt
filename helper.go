package tele_prompt

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler[T any] func() (T, error)

type SendMessageConfig[T any] struct {
	InitialMessage    tgbotapi.Chattable
	OnErrorMessage    tgbotapi.Chattable
	OnMaxRetryMessage tgbotapi.Chattable
	Handler           Handler[T]
	MaxRetry          int
}

func SendMessage[T any](bot *tgbotapi.BotAPI, config SendMessageConfig[T]) (T, error) {
	var t T
	_, err := bot.Send(config.InitialMessage)
	if err != nil {
		return t, err
	}
	retryCount := 0
	for {
		t, err := config.Handler()

		if err != nil {
			// if OnErrorMessage is not set, SendMessage will not handle errors
			if config.OnErrorMessage == nil {
				return t, err
			}

			// if OnErrorMessage is set, SendMessage will retry
			retryCount++
			if config.MaxRetry > 0 && retryCount > config.MaxRetry {
				if config.OnMaxRetryMessage != nil {
					bot.Send(config.OnMaxRetryMessage)
				}
				return t, ErrMaxRetries
			}
			if _, err := bot.Send(config.OnErrorMessage); err != nil {
				return t, err
			}
			continue
		}
		return t, nil
	}

}
