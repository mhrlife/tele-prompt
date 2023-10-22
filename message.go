package tele_prompt

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *Manager) TextField(ctx context.Context, userId int64) (string, error) {
	update, err := m.Any(ctx, userId)
	if err != nil {
		return "", err
	}
	return textField(update)
}

func textField(update tgbotapi.Update) (string, error) {
	if update.Message == nil || update.Message.Text == "" {
		return "", ErrBadResponse
	}

	return update.Message.Text, nil
}

// Image waits for an image response from the user and returns the file ID of the highest resolution image.
func (m *Manager) Image(ctx context.Context, userId int64) (tgbotapi.PhotoSize, error) {
	update, err := m.Any(ctx, userId)
	if err != nil {
		return tgbotapi.PhotoSize{}, err
	}
	return highestResolutionImage(update)
}

// highestResolutionImage parses the provided update to extract the highest resolution image's file ID.
func highestResolutionImage(update tgbotapi.Update) (tgbotapi.PhotoSize, error) {
	if update.Message == nil || update.Message.Photo == nil {
		return tgbotapi.PhotoSize{}, ErrBadResponse
	}

	// Assuming that the last photo in the array is the highest resolution.
	// The PhotoSize array is usually ordered by file size.
	highestRes := update.Message.Photo[len(update.Message.Photo)-1]
	return highestRes, nil
}

// Voice waits for a voice message response from the user and returns the file ID of the voice message.
func (m *Manager) Voice(ctx context.Context, userId int64) (*tgbotapi.Voice, error) {
	update, err := m.Any(ctx, userId)
	if err != nil {
		return nil, err
	}
	return extractVoice(update)
}

// extractVoiceFileID parses the provided update to extract the voice message's file ID.
func extractVoice(update tgbotapi.Update) (*tgbotapi.Voice, error) {
	if update.Message == nil || update.Message.Voice == nil {
		return nil, ErrBadResponse
	}

	return update.Message.Voice, nil
}

// Video waits for a video response from the user and returns the file ID of the video.
func (m *Manager) Video(ctx context.Context, userId int64) (*tgbotapi.Video, error) {
	update, err := m.Any(ctx, userId)
	if err != nil {
		return nil, err
	}
	return extractVideo(update)
}

func extractVideo(update tgbotapi.Update) (*tgbotapi.Video, error) {
	if update.Message == nil || update.Message.Video == nil {
		return nil, ErrBadResponse
	}
	return update.Message.Video, nil
}

// VideoNote waits for a video note response from the user and returns the file ID of the video note.
func (m *Manager) VideoNote(ctx context.Context, userId int64) (*tgbotapi.VideoNote, error) {
	update, err := m.Any(ctx, userId)
	if err != nil {
		return nil, err
	}
	return extractVideoNote(update)
}

func extractVideoNote(update tgbotapi.Update) (*tgbotapi.VideoNote, error) {
	if update.Message == nil || update.Message.VideoNote == nil {
		return nil, ErrBadResponse
	}
	return update.Message.VideoNote, nil
}

// Document waits for a document response from the user and returns the file ID of the document.
func (m *Manager) Document(ctx context.Context, userId int64) (*tgbotapi.Document, error) {
	update, err := m.Any(ctx, userId)
	if err != nil {
		return nil, err
	}
	return extractDocument(update)
}

func extractDocument(update tgbotapi.Update) (*tgbotapi.Document, error) {
	if update.Message == nil || update.Message.Document == nil {
		return nil, ErrBadResponse
	}
	return update.Message.Document, nil
}

func (m *Manager) Location(ctx context.Context, userId int64) (*tgbotapi.Location, error) {
	update, err := m.Any(ctx, userId)
	if err != nil {
		return nil, err
	}
	return extractLocation(update)
}

func extractLocation(update tgbotapi.Update) (*tgbotapi.Location, error) {
	if update.Message == nil || update.Message.Location == nil {
		return nil, ErrValidationFailed
	}
	return update.Message.Location, nil
}
