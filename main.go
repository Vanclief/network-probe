package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)

	app := cli.NewApp()
	app.Name = "probius"
	app.Usage = "Network prober"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "destination",
			Aliases:  []string{"d"},
			Usage:    "Address to probe",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "influx-url",
			Aliases: []string{"influx"},
			Usage:   "Enable sending logs to influx",
		},
		&cli.StringFlag{
			Name:    "influx-bucket",
			Aliases: []string{"bucket"},
			Usage:   "The bucket that will receive the logs",
		},
	}

	app.Action = func(c *cli.Context) error {
		probe, err := NewProbe()
		if err != nil {
			return cli.Exit(err, 1)
		}

		err = probe.Run(c.String("destination"))
		if err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(0)
	}
}
