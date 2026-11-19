package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// maxDestinationLength is the longest accepted destination URL, in bytes.
const maxDestinationLength = 2048

// Validation errors returned by ValidateDestination. Handlers map these to
// client-facing messages, so they must stay distinguishable.
var (
	ErrEmptyDestination   = errors.New("destination URL is empty")
	ErrDestinationTooLong = errors.New("destination URL is longer than 2048 bytes")
	ErrInvalidDestination = errors.New("destination URL is not a valid absolute URL")
	ErrUnsupportedScheme  = errors.New("destination URL must use http or https")
	ErrMissingHost        = errors.New("destination URL is missing a host")
)

// ShortURL is the stored mapping between a short code and the destination URL
// it redirects to.
type ShortURL struct {
	// Code is the generated identifier that appears in the short URL path.
	Code string
	// OriginalURL is the validated destination the short URL redirects to.
	OriginalURL string
}

// ValidateDestination reports whether raw is an acceptable redirect target and
// returns the trimmed value to store. A destination must be an absolute http or
// https URL with a non-empty host, and must not exceed maxDestinationLength
// bytes. The submitted URL is never fetched or inspected beyond parsing.
func ValidateDestination(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", ErrEmptyDestination
	}
	if len(trimmed) > maxDestinationLength {
		return "", ErrDestinationTooLong
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidDestination, err)
	}
	if !parsed.IsAbs() {
		return "", fmt.Errorf("%w: missing scheme", ErrInvalidDestination)
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("%w: got %q", ErrUnsupportedScheme, parsed.Scheme)
	}
	if parsed.Host == "" {
		return "", ErrMissingHost
	}

	return trimmed, nil
}
