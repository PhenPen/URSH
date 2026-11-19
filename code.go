package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// codeLength is the number of characters in every generated short code.
const codeLength = 7

// codeAlphabet is the character set short codes are drawn from, base62.
// Letters and digits need no escaping in a URL path. With 7 characters the
// space is 62^7, about 3.5 trillion codes, so collisions stay rare.
const codeAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// NewCode returns a random short code made of characters from codeAlphabet.
// The randomness comes from crypto/rand, not math/rand, because guessed codes
// would let anyone enumerate other people's links.
func NewCode() (string, error) {
	return newCode(rand.Reader)
}

// newCode is NewCode with an injectable random source. Tests supply either
// fixed bytes or a reader that fails, so they do not depend on real randomness.
func newCode(random io.Reader) (string, error) {
	limit := big.NewInt(int64(len(codeAlphabet)))
	code := make([]byte, codeLength)
	for i := range code {
		// rand.Int returns a uniform value in [0, 62) without the modulo bias
		// that byte % 62 would introduce.
		n, err := rand.Int(random, limit)
		if err != nil {
			return "", fmt.Errorf("generating short code: %w", err)
		}
		code[i] = codeAlphabet[n.Int64()]
	}
	return string(code), nil
}
