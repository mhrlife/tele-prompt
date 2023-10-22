package tele_prompt

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestManager_Any(t *testing.T) {
	m := NewManager(DefaultConfig())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		update, err := m.Any(context.Background(), 12)
		assert.NoError(t, err)
		assert.Equal(t, "test", update.Message.Text)
		wg.Done()
	}()

	go func() {
		<-time.After(time.Millisecond * 100)
		// this update is from another user and should be skipped
		hf := m.Handle(tgbotapi.Update{
			Message: &tgbotapi.Message{
				From: &tgbotapi.User{
					ID: 13,
				},
				Text: "test",
			},
		})

		// this update has to be processed
		assert.Equal(t, false, hf)
		h1 := m.Handle(tgbotapi.Update{
			Message: &tgbotapi.Message{
				From: &tgbotapi.User{
					ID: 12,
				},
				Text: "test",
			},
		})
		assert.Equal(t, h1, true)
		wg.Done()
	}()
	wg.Wait()
}

func TestManager_Any_Timeout(t *testing.T) {
	m := NewManager(Config{DefaultTimeout: time.Millisecond})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, err := m.Any(context.Background(), 12)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		wg.Done()
	}()
	wg.Wait()
}
