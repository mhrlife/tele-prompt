package tele_prompt

import (
	"context"
	"net/url"
	"regexp"
	"strconv"
)

func (m *Manager) OneOf(ctx context.Context, userId int64, options []string) (string, error) {
	txt, err := m.TextField(ctx, userId)
	if err != nil {
		return "", err
	}
	return oneOf(txt, options)
}

func oneOf(txt string, options []string) (string, error) {
	if !in(txt, options) {
		return "", ErrValidationFailed
	}

	return txt, nil
}

func in[T comparable](item T, arr []T) bool {
	for _, t := range arr {
		if t == item {
			return true
		}
	}
	return false
}

func (m *Manager) Range(ctx context.Context, userId int64, min, max int) (int, error) {
	txt, err := m.TextField(ctx, userId)
	if err != nil {
		return 0, err
	}
	return inRange(txt, min, max)
}

func inRange(txt string, min, max int) (int, error) {
	num, err := strconv.Atoi(txt)
	if err != nil {
		return 0, ErrValidationFailed
	}

	if num < min || num > max {
		return 0, ErrOutofRange
	}

	return num, nil
}

// Url prompts the user for a URL and validates it.
func (m *Manager) Url(ctx context.Context, userId int64) (string, error) {
	txt, err := m.TextField(ctx, userId)
	if err != nil {
		return "", err
	}
	return validateUrl(txt)
}

func validateUrl(txt string) (string, error) {
	_, err := url.ParseRequestURI(txt)
	if err != nil {
		return "", ErrValidationFailed
	}
	return txt, nil
}

// Email prompts the user for an email and validates it.
func (m *Manager) Email(ctx context.Context, userId int64) (string, error) {
	txt, err := m.TextField(ctx, userId)
	if err != nil {
		return "", err
	}
	return validateEmail(txt)
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func validateEmail(txt string) (string, error) {
	// Using regex to match common email format
	if !emailRegex.MatchString(txt) {
		return "", ErrValidationFailed
	}
	return txt, nil
}
