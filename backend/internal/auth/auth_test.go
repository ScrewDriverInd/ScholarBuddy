package auth

import "testing"

func TestPrincipalHasRole(t *testing.T) {
	p := Principal{Roles: []string{"ROLE_USER", "ROLE_ADMIN"}}
	if !p.HasRole("ROLE_ADMIN") || !p.HasRole("ROLE_USER") || p.HasRole("ROLE_UNKNOWN") {
		t.Fatal("role membership is incorrect")
	}
}
