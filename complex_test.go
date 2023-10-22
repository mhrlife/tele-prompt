package tele_prompt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOneOf(t *testing.T) {
	_, err := oneOf("x", []string{"y", "z"})
	assert.ErrorIs(t, ErrValidationFailed, err)

	str, err := oneOf("x", []string{"x", "y", "z"})
	assert.NoError(t, err)
	assert.Equal(t, "x", str)
}

func TestInRange(t *testing.T) {
	_, err := inRange("5", 6, 10)
	assert.ErrorIs(t, ErrOutofRange, err)

	_, err = inRange("abc", 1, 10)
	assert.ErrorIs(t, ErrValidationFailed, err)

	num, err := inRange("7", 5, 10)
	assert.NoError(t, err)
	assert.Equal(t, 7, num)
}

func TestValidateUrl(t *testing.T) {
	_, err := validateUrl("not_a_url")
	assert.ErrorIs(t, ErrValidationFailed, err)

	urlStr, err := validateUrl("https://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, "https://www.google.com", urlStr)
}

func TestValidateEmail(t *testing.T) {
	_, err := validateEmail("not_an_email")
	assert.ErrorIs(t, ErrValidationFailed, err)

	emailStr, err := validateEmail("example@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "example@example.com", emailStr)
}
