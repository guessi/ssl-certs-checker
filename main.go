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
				Name:     "hosts",
				Aliases:  []string{"H"},
				Value:    "",
				Usage:    "target hosts, splits by comma",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			prettyPrintCertsInfo(c.String("hosts"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
