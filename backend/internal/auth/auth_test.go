package auth

import "testing"

func TestCredentialsMatch(t *testing.T) {
	if !CredentialsMatch("admin", "password", "key", "admin", "password", "key") {
		t.Fatal("expected valid credentials to match")
	}
	if CredentialsMatch("admin", "password", "key", "admin", "wrong", "key") {
		t.Fatal("expected invalid credentials not to match")
	}
}
