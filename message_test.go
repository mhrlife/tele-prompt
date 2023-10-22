package tele_prompt

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextField(t *testing.T) {
	_, err := textField(tgbotapi.Update{})
	assert.ErrorIs(t, ErrBadResponse, err)

	str, err := textField(tgbotapi.Update{Message: &tgbotapi.Message{Text: "hello"}})
	assert.NoError(t, err)
	assert.Equal(t, "hello", str)
}

// Test for highestResolutionImage
func TestHighestResolutionImage(t *testing.T) {
	_, err := highestResolutionImage(tgbotapi.Update{})
	assert.ErrorIs(t, ErrBadResponse, err)

	photoSize := tgbotapi.PhotoSize{
		FileID:   "testFileId",
		Width:    1280,
		Height:   720,
		FileSize: 50000,
	}

	item, err := highestResolutionImage(tgbotapi.Update{Message: &tgbotapi.Message{Photo: []tgbotapi.PhotoSize{photoSize}}})
	assert.NoError(t, err)
	assert.Equal(t, "testFileId", item.FileID)
}

// Test for extractVoiceFileID
func TestExtractVoice(t *testing.T) {
	_, err := extractVoice(tgbotapi.Update{})
	assert.ErrorIs(t, ErrBadResponse, err)

	voice := tgbotapi.Voice{
		FileID:   "testVoiceFileId",
		Duration: 30,
		FileSize: 15000,
	}

	item, err := extractVoice(tgbotapi.Update{Message: &tgbotapi.Message{Voice: &voice}})
	assert.NoError(t, err)
	assert.Equal(t, "testVoiceFileId", item.FileID)
}

func TestExtractVideo(t *testing.T) {
	_, err := extractVideo(tgbotapi.Update{})
	assert.ErrorIs(t, ErrBadResponse, err)

	video := tgbotapi.Video{
		FileID: "testVideoFileId",
		Width:  1280,
		Height: 720,
	}

	item, err := extractVideo(tgbotapi.Update{Message: &tgbotapi.Message{Video: &video}})
	assert.NoError(t, err)
	assert.Equal(t, "testVideoFileId", item.FileID)
}

func TestExtractVideoNote(t *testing.T) {
	_, err := extractVideoNote(tgbotapi.Update{})
	assert.ErrorIs(t, ErrBadResponse, err)

	videoNote := tgbotapi.VideoNote{
		FileID: "testVideoNoteFileId",
		Length: 360,
	}

	item, err := extractVideoNote(tgbotapi.Update{Message: &tgbotapi.Message{VideoNote: &videoNote}})
	assert.NoError(t, err)
	assert.Equal(t, "testVideoNoteFileId", item.FileID)
}

func TestExtractDocument(t *testing.T) {
	_, err := extractDocument(tgbotapi.Update{})
	assert.ErrorIs(t, ErrBadResponse, err)

	document := tgbotapi.Document{
		FileID: "testDocumentFileId",
	}

	doc, err := extractDocument(tgbotapi.Update{Message: &tgbotapi.Message{Document: &document}})
	assert.NoError(t, err)
	assert.Equal(t, "testDocumentFileId", doc.FileID)
}
