package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/guessi/ssl-certs-checker/pkg/config"
)

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}

	if app.formatter == nil {
		t.Error("New() should initialize formatter")
	}
}

func TestApp_Run_InvalidConfig(t *testing.T) {
	app := New()
	ctx := context.Background()

	// Test with invalid config (no domains or config file)
	cfg := &config.AppConfig{
		Timeout:      5,
		OutputFormat: "table",
	}

	err := app.Run(ctx, cfg)
	if err == nil {
		t.Error("Run() should return error for invalid config")
	}
}

func TestApp_Run_ValidConfigWithDomains(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	app := New()
	ctx := context.Background()

	// Test with valid config using domains
	cfg := &config.AppConfig{
		Domains:      "google.com",
		Timeout:      10,
		OutputFormat: "json",
	}

	err := app.Run(ctx, cfg)
	if err != nil {
		t.Errorf("Run() unexpected error with valid config: %v", err)
	}
}

func TestApp_Run_ValidConfigWithFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "ssl-cert-app-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test.yaml")
	configContent := `hosts:
  - google.com
  - github.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	app := New()
	ctx := context.Background()

	// Test with valid config using config file
	cfg := &config.AppConfig{
		ConfigFile:   configPath,
		Timeout:      10,
		OutputFormat: "table",
	}

	err = app.Run(ctx, cfg)
	if err != nil {
		t.Errorf("Run() unexpected error with valid config file: %v", err)
	}
}

func TestApp_Run_NonExistentConfigFile(t *testing.T) {
	app := New()
	ctx := context.Background()

	// Test with non-existent config file
	cfg := &config.AppConfig{
		ConfigFile:   "/non/existent/file.yaml",
		Timeout:      5,
		OutputFormat: "table",
	}

	err := app.Run(ctx, cfg)
	if err == nil {
		t.Error("Run() should return error for non-existent config file")
	}
}

func TestApp_Run_InvalidDomains(t *testing.T) {
	app := New()
	ctx := context.Background()

	// Test with invalid domains - the app should not error but should show errors in output
	cfg := &config.AppConfig{
		Domains:      "invalid::domain,another::invalid",
		Timeout:      5,
		OutputFormat: "table",
	}

	// The app should not return an error - it should handle invalid domains gracefully
	err := app.Run(ctx, cfg)
	if err != nil {
		t.Errorf("Run() should not return error for invalid domains, should handle gracefully: %v", err)
	}
}

func TestApp_Run_ContextCancellation(t *testing.T) {
	app := New()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	cfg := &config.AppConfig{
		Domains:      "google.com",
		Timeout:      30, // Long timeout
		OutputFormat: "table",
	}

	err := app.Run(ctx, cfg)
	if err == nil {
		t.Error("Run() should return error for cancelled context")
	}
}
