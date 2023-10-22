package tele_prompt

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrBadResponse      = errors.New("response is incompatible with the prompt")
	ErrOutofRange       = errors.New("response is out of range")
	ErrValidationFailed = errors.New("validation failed")
	ErrMaxRetries       = errors.New("max retries has been reached")
)

type Config struct {
	DefaultTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		DefaultTimeout: time.Minute * 5,
	}
}

type Manager struct {
	config   Config
	channels sync.Map
}

func NewManager(config Config) *Manager {
	return &Manager{
		config:   config,
		channels: sync.Map{},
	}
}

// Handle receives a new update and returns true if app expected user to answer to something
// Otherwise, returns false
func (m *Manager) Handle(update tgbotapi.Update) bool {
	userId := userIdFromUpdate(update)
	if userId == -1 {
		return false
	}

	channelInterface, has := m.channels.Load(userId)
	if !has {
		return false
	}

	select {
	case channelInterface.(chan tgbotapi.Update) <- update:
		return true
	default:
		log.Printf("teleprompt: channel was full, unexpected error\n")
		return false
	}

}

// Any does not care about types and different prompts, it receives the entire update
func (m *Manager) Any(ctx context.Context, userId int64) (tgbotapi.Update, error) {
	ctx, cancel := context.WithTimeout(ctx, m.config.DefaultTimeout)
	defer cancel()
	ch := make(chan tgbotapi.Update, 1)
	m.channels.Store(userId, ch)
	// delete the current channel
	defer func() {
		m.channels.Delete(userId)
	}()
	select {
	case <-ctx.Done():
		return tgbotapi.Update{}, ctx.Err()
	case update := <-ch:
		close(ch)
		return update, nil
	}
}

func userIdFromUpdate(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return -1
}
