package main

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateDestination(t *testing.T) {
	prefix := "https://example.com/"
	atLimit := prefix + strings.Repeat("a", maxDestinationLength-len(prefix))
	overLimit := prefix + strings.Repeat("a", maxDestinationLength-len(prefix)+1)

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "https url",
			input: "https://example.com/products/a-long-product-path",
			want:  "https://example.com/products/a-long-product-path",
		},
		{
			name:  "http url",
			input: "http://example.com",
			want:  "http://example.com",
		},
		{
			name:  "uppercase scheme",
			input: "HTTPS://example.com",
			want:  "HTTPS://example.com",
		},
		{
			name:  "query and fragment",
			input: "https://example.com/search?q=go#results",
			want:  "https://example.com/search?q=go#results",
		},
		{
			name:  "surrounding spaces are trimmed",
			input: "  https://example.com  ",
			want:  "https://example.com",
		},
		{
			name:  "host with port",
			input: "http://localhost:8080/health",
			want:  "http://localhost:8080/health",
		},
		{
			name:  "exactly at the length limit",
			input: atLimit,
			want:  atLimit,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: ErrEmptyDestination,
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: ErrEmptyDestination,
		},
		{
			name:    "over the length limit",
			input:   overLimit,
			wantErr: ErrDestinationTooLong,
		},
		{
			name:    "relative path",
			input:   "/products/1",
			wantErr: ErrInvalidDestination,
		},
		{
			name:    "missing scheme",
			input:   "example.com",
			wantErr: ErrInvalidDestination,
		},
		{
			name:    "invalid escape",
			input:   "http://example.com/%zz",
			wantErr: ErrInvalidDestination,
		},
		{
			name:    "ftp scheme",
			input:   "ftp://example.com/file",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "javascript scheme",
			input:   "javascript:alert(1)",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "missing host",
			input:   "http://",
			wantErr: ErrMissingHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateDestination(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ValidateDestination(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				if got != "" {
					t.Errorf("ValidateDestination(%q) returned %q with an error, want empty string", tt.input, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("ValidateDestination(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ValidateDestination(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
