package config

import "testing"

func TestProductionRejectsDefaultSecrets(t *testing.T) {
	cfg := &Config{
		Environment:      "production",
		JWTSecret:        "your-default-secret-change-in-production",
		JWTRefreshSecret: "your-refresh-secret-change-in-production",
		DatabaseURL:      "postgres://example",
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("production accepted default JWT secrets")
	}
}

func TestProductionAcceptsCompleteConfiguration(t *testing.T) {
	cfg := &Config{
		Environment:      "production",
		JWTSecret:        "a-secure-access-secret-with-32-characters",
		JWTRefreshSecret: "a-secure-refresh-secret-with-32-characters",
		DatabaseURL:      "postgres://example",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid production config rejected: %v", err)
	}
}
