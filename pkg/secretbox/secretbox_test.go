package secretbox

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

func testKey(seed byte) string {
	key := make([]byte, KeySizeAES256GCM)
	for index := range key {
		key[index] = seed
	}
	return base64.StdEncoding.EncodeToString(key)
}

func TestAES256GCMSecretBoxRoundTrip(t *testing.T) {
	key, err := DecodeBase64Key(testKey(7))
	if err != nil {
		t.Fatalf("decode key: %v", err)
	}
	box, err := NewAES256GCM(key, "totp-test")
	if err != nil {
		t.Fatalf("create box: %v", err)
	}

	envelope, err := box.EncryptString("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if !strings.HasPrefix(envelope, "nanami-secretbox-v1:totp-test:") {
		t.Fatalf("expected versioned envelope, got %q", envelope)
	}
	if strings.Contains(envelope, "JBSWY3DPEHPK3PXP") {
		t.Fatalf("ciphertext envelope must not contain plaintext secret: %q", envelope)
	}
	plaintext, err := box.DecryptString(envelope)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if plaintext != "JBSWY3DPEHPK3PXP" {
		t.Fatalf("unexpected plaintext %q", plaintext)
	}
}

func TestAES256GCMSecretBoxRejectsPlaintextAndWrongKey(t *testing.T) {
	keyA, err := DecodeBase64Key(testKey(1))
	if err != nil {
		t.Fatalf("decode key A: %v", err)
	}
	keyB, err := DecodeBase64Key(testKey(2))
	if err != nil {
		t.Fatalf("decode key B: %v", err)
	}
	boxA, err := NewAES256GCM(keyA, "totp-test")
	if err != nil {
		t.Fatalf("create box A: %v", err)
	}
	boxB, err := NewAES256GCM(keyB, "totp-test")
	if err != nil {
		t.Fatalf("create box B: %v", err)
	}

	envelope, err := boxA.EncryptString("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if _, err := boxA.DecryptString("JBSWY3DPEHPK3PXP"); !errors.Is(err, ErrInvalidEnvelope) {
		t.Fatalf("expected plaintext rejection, got %v", err)
	}
	if _, err := boxB.DecryptString(envelope); !errors.Is(err, ErrInvalidEnvelope) {
		t.Fatalf("expected wrong key rejection, got %v", err)
	}
}

func TestDecodeBase64KeyRequiresAES256Material(t *testing.T) {
	if _, err := DecodeBase64Key(""); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("expected empty key rejection, got %v", err)
	}
	if _, err := DecodeBase64Key(base64.StdEncoding.EncodeToString([]byte("short"))); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("expected short key rejection, got %v", err)
	}
}
