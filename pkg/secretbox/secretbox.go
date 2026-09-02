package secretbox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	KeySizeAES256GCM = 32

	envelopeVersion = "nanami-secretbox-v1"
)

var (
	ErrInvalidKey      = errors.New("invalid secretbox key")
	ErrInvalidEnvelope = errors.New("invalid secretbox envelope")

	keyIDPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,64}$`)
)

// Box encrypts short server-side secret material into versioned AES-GCM envelopes.
type Box struct {
	aead  cipher.AEAD
	keyID string
	rand  io.Reader
}

// DecodeBase64Key decodes an operator-managed AES-256 key.
func DecodeBase64Key(raw string) ([]byte, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, ErrInvalidKey
	}
	key, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("%w: key must be standard base64", ErrInvalidKey)
	}
	if len(key) != KeySizeAES256GCM {
		return nil, fmt.Errorf("%w: decoded key must be %d bytes", ErrInvalidKey, KeySizeAES256GCM)
	}
	return key, nil
}

// NewAES256GCM creates an AES-256-GCM secret box.
func NewAES256GCM(key []byte, keyID string) (*Box, error) {
	if len(key) != KeySizeAES256GCM {
		return nil, fmt.Errorf("%w: decoded key must be %d bytes", ErrInvalidKey, KeySizeAES256GCM)
	}
	keyID = strings.TrimSpace(keyID)
	if keyID == "" {
		keyID = "default"
	}
	if !keyIDPattern.MatchString(keyID) {
		return nil, fmt.Errorf("%w: key id must match %s", ErrInvalidKey, keyIDPattern.String())
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Box{aead: aead, keyID: keyID, rand: rand.Reader}, nil
}

// EncryptString encrypts plaintext into a self-describing envelope.
func (b *Box) EncryptString(plaintext string) (string, error) {
	if b == nil || b.aead == nil {
		return "", ErrInvalidKey
	}
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := io.ReadFull(b.rand, nonce); err != nil {
		return "", err
	}
	header := envelopeVersion + ":" + b.keyID
	ciphertext := b.aead.Seal(nil, nonce, []byte(plaintext), []byte(header))
	return header + ":" +
		base64.RawURLEncoding.EncodeToString(nonce) + ":" +
		base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// DecryptString decrypts a valid envelope. Non-envelope plaintext is rejected.
func (b *Box) DecryptString(envelope string) (string, error) {
	if b == nil || b.aead == nil {
		return "", ErrInvalidKey
	}
	parts := strings.Split(strings.TrimSpace(envelope), ":")
	if len(parts) != 4 || parts[0] != envelopeVersion || parts[1] == "" {
		return "", ErrInvalidEnvelope
	}
	if parts[1] != b.keyID {
		return "", fmt.Errorf("%w: unexpected key id", ErrInvalidEnvelope)
	}
	nonce, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "", fmt.Errorf("%w: invalid nonce", ErrInvalidEnvelope)
	}
	if len(nonce) != b.aead.NonceSize() {
		return "", fmt.Errorf("%w: invalid nonce size", ErrInvalidEnvelope)
	}
	ciphertext, err := base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil {
		return "", fmt.Errorf("%w: invalid ciphertext", ErrInvalidEnvelope)
	}
	header := parts[0] + ":" + parts[1]
	plaintext, err := b.aead.Open(nil, nonce, ciphertext, []byte(header))
	if err != nil {
		return "", fmt.Errorf("%w: authentication failed", ErrInvalidEnvelope)
	}
	return string(plaintext), nil
}

// IsEncrypted reports whether a value has the current envelope shape.
func (b *Box) IsEncrypted(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), envelopeVersion+":")
}
