package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	// Alphabet defines the characters used in the short URL code
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// CodeLength defines the length of the generated short code
	CodeLength = 7
)

// GenerateShortCode creates a random short code for URLs
func GenerateShortCode() (string, error) {
	code := strings.Builder{}
	alphabetLen := big.NewInt(int64(len(Alphabet)))

	for i := 0; i < CodeLength; i++ {
		// Generate a random index within the alphabet
		randomIndex, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}

		// Add the character at the random index to our code
		code.WriteByte(Alphabet[randomIndex.Int64()])
	}

	return code.String(), nil
}

// ValidateURL performs simple validation on a URL
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
