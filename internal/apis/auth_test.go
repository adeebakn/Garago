package apis

import (
	"strings"
	"testing"
)

func TestPasswordResetTokenKeyDoesNotContainToken(t *testing.T) {
	token := "a-sensitive-reset-token"
	if key := passwordResetTokenKey(token); strings.Contains(key, token) {
		t.Fatal("reset token key must not contain the plaintext token")
	}
}
