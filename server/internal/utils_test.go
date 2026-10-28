package internal

import "testing"

func TestEncryptStringRoundTrip(t *testing.T) {
	plaintext := "sensitive value"
	password := "test password"

	encrypted, err := EncryptString(plaintext, password)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	if encrypted == plaintext {
		t.Fatal("EncryptString() returned the plaintext")
	}

	decrypted, err := DecryptString(encrypted, password)
	if err != nil {
		t.Fatalf("DecryptString() error = %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("DecryptString() = %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptStringRejectsWrongPasswordAndTampering(t *testing.T) {
	encrypted, err := EncryptString("sensitive value", "correct password")
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}

	if _, err := DecryptString(encrypted, "wrong password"); err == nil {
		t.Error("DecryptString() accepted the wrong password")
	}

	tampered := []byte(encrypted)
	if tampered[len(encryptedStringPrefix)] == 'A' {
		tampered[len(encryptedStringPrefix)] = 'B'
	} else {
		tampered[len(encryptedStringPrefix)] = 'A'
	}
	if _, err := DecryptString(string(tampered), "correct password"); err == nil {
		t.Error("DecryptString() accepted a tampered value")
	}
}
