// Package encrypt provides simple symmetric encryption and decryption
// for secret values stored in the secret-store backend.
//
// It uses AES-256-GCM with a key derived from a passphrase via PBKDF2.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// saltSize is the number of random bytes used as the PBKDF2 salt.
	saltSize = 16
	// keySize is the AES-256 key length in bytes.
	keySize = 32
	// pbkdf2Iter is the number of PBKDF2 iterations.
	pbkdf2Iter = 100_000
)

// ErrInvalidCiphertext is returned when decryption fails due to malformed input.
var ErrInvalidCiphertext = errors.New("encrypt: invalid ciphertext")

// deriveKey derives a 32-byte AES key from passphrase and salt using PBKDF2-SHA256.
func deriveKey(passphrase, salt []byte) []byte {
	return pbkdf2.Key(passphrase, salt, pbkdf2Iter, keySize, sha256.New)
}

// Encrypt encrypts plaintext using AES-256-GCM with a key derived from passphrase.
// The returned string is base64-encoded and contains the salt, nonce, and ciphertext.
func Encrypt(passphrase, plaintext string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := deriveKey([]byte(passphrase), salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	payload := append(salt, ciphertext...) //nolint:gocritic
	return base64.StdEncoding.EncodeToString(payload), nil
}

// Decrypt decrypts a base64-encoded ciphertext produced by Encrypt.
func Decrypt(passphrase, encoded string) (string, error) {
	payload, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	if len(payload) < saltSize {
		return "", ErrInvalidCiphertext
	}

	salt, data := payload[:saltSize], payload[saltSize:]
	key := deriveKey([]byte(passphrase), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", ErrInvalidCiphertext
	}

	plaintext, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plaintext), nil
}
