package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "SSL Certificate Checker",
		Usage: "check SSL certificates at once",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"C"},
				Value:    "",
				Usage:    "config file",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "timeout",
				Aliases:  []string{"t"},
				Value:    defaultDialerTimeout,
				Usage:    "dialer timeout in second(s)",
				Required: false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			prettyPrintCertsInfo(c.String("config"), c.Int("timeout"))
			return nil
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		os.Exit(1)
	}
}
