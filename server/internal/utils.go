package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	encryptedStringPrefix = "v1:"
	argon2Time            = 3
	argon2Memory          = 64 * 1024
	argon2Threads         = 1
	argon2KeyLength       = 32
	argon2SaltLength      = 16
)

func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Map[T, U any](in []T, f func(T) U) []U {
	out := make([]U, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func HashString(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:]) // 64 chars
}

func EncryptString(plaintext, password string) (string, error) {
	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate encryption salt: %w", err)
	}

	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLength)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create encryption cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create encryption mode: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generate encryption nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), []byte(encryptedStringPrefix))
	payload := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	payload = append(payload, salt...)
	payload = append(payload, nonce...)
	payload = append(payload, ciphertext...)

	return encryptedStringPrefix + base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecryptString(encrypted, password string) (string, error) {
	if !strings.HasPrefix(encrypted, encryptedStringPrefix) {
		return "", errors.New("unsupported encrypted string format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(encrypted, encryptedStringPrefix))
	if err != nil {
		return "", fmt.Errorf("decode encrypted string: %w", err)
	}

	minimumLength := argon2SaltLength + 12 + 16
	if len(payload) < minimumLength {
		return "", errors.New("encrypted string is too short")
	}
	salt := payload[:argon2SaltLength]
	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLength)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create decryption cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create decryption mode: %w", err)
	}

	nonceEnd := argon2SaltLength + gcm.NonceSize()
	if len(payload) < nonceEnd+gcm.Overhead() {
		return "", errors.New("encrypted string is too short")
	}
	plaintext, err := gcm.Open(nil, payload[argon2SaltLength:nonceEnd], payload[nonceEnd:], []byte(encryptedStringPrefix))
	if err != nil {
		return "", errors.New("decrypt encrypted string: authentication failed")
	}
	return string(plaintext), nil
}
