package app

import (
	"context"
	"fmt"
	"time"

	"github.com/guessi/ssl-certs-checker/pkg/cert"
	"github.com/guessi/ssl-certs-checker/pkg/config"
	"github.com/guessi/ssl-certs-checker/pkg/output"
)

// New creates a new application instance
func New() *App {
	return &App{
		formatter: output.New(),
	}
}

// Run executes the application with the given configuration
func (a *App) Run(ctx context.Context, cfg *config.AppConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	hosts, err := cfg.GetHosts()
	if err != nil {
		return fmt.Errorf("failed to get hosts: %w", err)
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	a.checker = cert.New(timeout, cfg.Insecure)

	result, err := a.checker.CheckCertificates(ctx, hosts)
	if err != nil {
		return fmt.Errorf("failed to check certificates: %w", err)
	}

	if err := a.formatter.Format(result, cfg.OutputFormat); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
