package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalToken := os.Getenv("TELEGRAM_TOKEN")
	originalDBPath := os.Getenv("DATABASE_PATH")

	// Clean up after test
	defer func() {
		os.Setenv("TELEGRAM_TOKEN", originalToken)
		os.Setenv("DATABASE_PATH", originalDBPath)
	}()

	// Test with environment variables
	os.Setenv("TELEGRAM_TOKEN", "test_token_123")
	os.Setenv("DATABASE_PATH", "/tmp/test.db")

	cfg := Load()

	if cfg.TelegramToken != "test_token_123" {
		t.Errorf("Expected token 'test_token_123', got '%s'", cfg.TelegramToken)
	}

	if cfg.DatabasePath != "/tmp/test.db" {
		t.Errorf("Expected database path '/tmp/test.db', got '%s'", cfg.DatabasePath)
	}

	if cfg.MaxResults != 3 {
		t.Errorf("Expected MaxResults 3, got %d", cfg.MaxResults)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("TELEGRAM_TOKEN")
	os.Unsetenv("DATABASE_PATH")
	os.Unsetenv("EMBEDDING_API_URL")
	os.Unsetenv("EMBEDDING_MODEL")

	cfg := Load()

	// Check defaults
	if cfg.TelegramToken != "" {
		t.Errorf("Expected empty token, got '%s'", cfg.TelegramToken)
	}

	if cfg.DatabasePath != "./messages.db" {
		t.Errorf("Expected default database path './messages.db', got '%s'", cfg.DatabasePath)
	}

	if cfg.EmbeddingAPIURL != "http://localhost:11434" {
		t.Errorf("Expected default API URL, got '%s'", cfg.EmbeddingAPIURL)
	}

	if cfg.EmbeddingModel != "all-minilm:latest" {
		t.Errorf("Expected default model, got '%s'", cfg.EmbeddingModel)
	}

	if cfg.MaxResults != 3 {
		t.Errorf("Expected MaxResults 3, got %d", cfg.MaxResults)
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with non-existing env var
	result = getEnv("NON_EXISTING_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}

	// Clean up
	os.Unsetenv("TEST_VAR")
}
