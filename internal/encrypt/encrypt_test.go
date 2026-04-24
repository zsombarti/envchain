package encrypt_test

import (
	"strings"
	"testing"

	"envchain/internal/encrypt"
)

const testPassphrase = "super-secret-passphrase"

func TestEncryptDecryptRoundtrip(t *testing.T) {
	plaintext := "MY_SECRET_VALUE"

	encoded, err := encrypt.Encrypt(testPassphrase, plaintext)
	if err != nil {
		t.Fatalf("Encrypt: unexpected error: %v", err)
	}

	got, err := encrypt.Decrypt(testPassphrase, encoded)
	if err != nil {
		t.Fatalf("Decrypt: unexpected error: %v", err)
	}
	if got != plaintext {
		t.Errorf("Decrypt = %q, want %q", got, plaintext)
	}
}

func TestEncryptProducesUniqueOutputs(t *testing.T) {
	// Each call should produce a different ciphertext due to random salt/nonce.
	a, err := encrypt.Encrypt(testPassphrase, "value")
	if err != nil {
		t.Fatal(err)
	}
	b, err := encrypt.Encrypt(testPassphrase, "value")
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Error("expected unique ciphertexts for identical inputs")
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	encoded, err := encrypt.Encrypt(testPassphrase, "secret")
	if err != nil {
		t.Fatal(err)
	}

	_, err = encrypt.Decrypt("wrong-passphrase", encoded)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	_, err := encrypt.Decrypt(testPassphrase, "!!!not-base64!!!")
	if err != encrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecryptTruncatedPayload(t *testing.T) {
	_, err := encrypt.Decrypt(testPassphrase, "dG9vc2hvcnQ=") // "tooshort" in base64
	if err != encrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestEncryptEmptyString(t *testing.T) {
	encoded, err := encrypt.Encrypt(testPassphrase, "")
	if err != nil {
		t.Fatalf("Encrypt empty: %v", err)
	}
	got, err := encrypt.Decrypt(testPassphrase, encoded)
	if err != nil {
		t.Fatalf("Decrypt empty: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestEncryptOutputIsBase64(t *testing.T) {
	encoded, err := encrypt.Encrypt(testPassphrase, "hello")
	if err != nil {
		t.Fatal(err)
	}
	// base64 standard encoding uses only A-Z a-z 0-9 + / =
	for _, c := range encoded {
		if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=", c) {
			t.Errorf("unexpected character in base64 output: %q", c)
		}
	}
}
