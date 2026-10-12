package mailer

import (
	"net/mail"
	"strings"
	"testing"
)

func TestBuildPasswordResetEmail(t *testing.T) {
	resetURL := "https://app.example.com/reset-password?token=abc123"
	message := buildPasswordResetEmail(mail.Address{Name: "Garago", Address: "notify@example.com"}, "user@example.com", resetURL)
	for _, expected := range []string{"To: user@example.com", "Subject: Reset your Garago password", resetURL, "text/plain", "text/html"} {
		if !strings.Contains(message, expected) {
			t.Fatalf("email message does not contain %q", expected)
		}
	}
}

func TestSMTPPasswordRemovesGoogleAppPasswordSpaces(t *testing.T) {
	if got := smtpPassword("smtp.gmail.com", "abcd efgh ijkl mnop"); got != "abcdefghijklmnop" {
		t.Fatalf("got %q", got)
	}
	if got := smtpPassword("smtp.example.com", "a password with spaces"); got != "a password with spaces" {
		t.Fatalf("non-Gmail password was modified: %q", got)
	}
}
