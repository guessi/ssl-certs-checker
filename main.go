package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
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
		Action: func(c *cli.Context) error {
			prettyPrintCertsInfo(c.String("config"), c.Int("timeout"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
