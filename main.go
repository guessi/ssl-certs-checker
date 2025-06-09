package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/guessi/ssl-certs-checker/pkg/app"
	"github.com/guessi/ssl-certs-checker/pkg/config"
	"github.com/urfave/cli/v3"
)

const defaultDialerTimeout = 5

func main() {
	cliApp := &cli.Command{
		Usage: "check SSL certificates at once",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"C"},
				Value:    "",
				Usage:    "config file",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "domains",
				Aliases:  []string{"d"},
				Value:    "",
				Usage:    "comma-separated list of domains to check (e.g., example.com,google.com:443)",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "timeout",
				Aliases:  []string{"t"},
				Value:    defaultDialerTimeout,
				Usage:    "dialer timeout in second(s)",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "insecure",
				Aliases:  []string{"k"},
				Value:    false,
				Usage:    "skip the verification of certificates",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Value:    "table",
				Usage:    "output format (table, json, yaml)",
				Required: false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			cfg := &config.AppConfig{
				ConfigFile:   c.String("config"),
				Domains:      c.String("domains"),
				Timeout:      c.Int("timeout"),
				Insecure:     c.Bool("insecure"),
				OutputFormat: c.String("output"),
			}

			// Create a context that can be cancelled by signals
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			// Handle graceful shutdown
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigChan
				cancel()
			}()

			application := app.New()
			if err := application.Run(ctx, cfg); err != nil {
				return cli.Exit(fmt.Sprintf("Error: %v", err), 1)
			}

			return nil
		},
	}

	if err := cliApp.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
