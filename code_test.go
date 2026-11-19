package main

import (
	"errors"
	"strings"
	"testing"
)

func TestCodeAlphabet(t *testing.T) {
	if got, want := len(codeAlphabet), 62; got != want {
		t.Errorf("len(codeAlphabet) = %d, want %d", got, want)
	}

	seen := make(map[rune]bool, len(codeAlphabet))
	for _, r := range codeAlphabet {
		if seen[r] {
			t.Errorf("codeAlphabet contains %q more than once", r)
		}
		seen[r] = true
	}
}

// TestNewCodeShape generates real codes and checks their properties. It never
// asserts a specific value, because the output is random by design.
func TestNewCodeShape(t *testing.T) {
	const runs = 200

	for i := 0; i < runs; i++ {
		code, err := NewCode()
		if err != nil {
			t.Fatalf("NewCode() error: %v", err)
		}
		if len(code) != codeLength {
			t.Fatalf("NewCode() = %q, length %d, want %d", code, len(code), codeLength)
		}
		for _, r := range code {
			if !strings.ContainsRune(codeAlphabet, r) {
				t.Fatalf("NewCode() = %q contains %q, which is outside codeAlphabet", code, r)
			}
		}
	}
}

func TestNewCodeVaries(t *testing.T) {
	first, err := NewCode()
	if err != nil {
		t.Fatalf("first NewCode() error: %v", err)
	}
	second, err := NewCode()
	if err != nil {
		t.Fatalf("second NewCode() error: %v", err)
	}

	if first == second {
		t.Errorf("two generated codes were equal (%q); with 62^7 possibilities that should be vanishingly rare", first)
	}
}

// repeatReader always fills the buffer with the same byte.
type repeatReader byte

func (r repeatReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(r)
	}
	return len(p), nil
}

// TestNewCodeDeterministic drives newCode with a fixed random source. Any byte
// below 62 is a valid alphabet index, so the mapping can be checked exactly.
func TestNewCodeDeterministic(t *testing.T) {
	tests := []struct {
		name  string
		value byte
		want  string
	}{
		{name: "lowest index", value: 0, want: "0000000"},
		{name: "highest index", value: 61, want: "zzzzzzz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newCode(repeatReader(tt.value))
			if err != nil {
				t.Fatalf("newCode() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("newCode() = %q, want %q", got, tt.want)
			}
		})
	}
}

// failingReader is an entropy source that always fails.
type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestNewCodeReportsEntropyFailure(t *testing.T) {
	wantErr := errors.New("entropy unavailable")

	code, err := newCode(failingReader{wantErr})

	if !errors.Is(err, wantErr) {
		t.Fatalf("newCode() error = %v, want an error wrapping %v", err, wantErr)
	}
	if code != "" {
		t.Errorf("newCode() returned %q with an error, want empty string", code)
	}
}
