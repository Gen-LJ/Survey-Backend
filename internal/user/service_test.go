package user

import (
	"strings"
	"testing"
)

// SeedAdmin validates its inputs before touching the database, so these paths
// are exercisable without one.
func TestSeedAdminNoopWhenUnconfigured(t *testing.T) {
	t.Setenv("ADMIN_EMAIL", "")
	t.Setenv("ADMIN_PASSWORD", "")

	if err := SeedAdmin(); err != nil {
		t.Fatalf("expected no-op, got %v", err)
	}
}

func TestSeedAdminNoopWhenOnlyEmailSet(t *testing.T) {
	t.Setenv("ADMIN_EMAIL", "admin@example.com")
	t.Setenv("ADMIN_PASSWORD", "")

	if err := SeedAdmin(); err != nil {
		t.Fatalf("expected no-op, got %v", err)
	}
}

func TestSeedAdminRejectsShortPassword(t *testing.T) {
	t.Setenv("ADMIN_EMAIL", "admin@example.com")
	t.Setenv("ADMIN_PASSWORD", "short")

	err := SeedAdmin()
	if err == nil {
		t.Fatal("expected a short password to be refused")
	}
	if !strings.Contains(err.Error(), "at least") {
		t.Fatalf("expected a length message, got %v", err)
	}
}
