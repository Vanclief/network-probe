package main

import (
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/urfave/cli/v2"
)

func main() {
	// Setup logs
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)

	// Setup env
	viper.AutomaticEnv()

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
			Name:     "server",
			Aliases:  []string{"s"},
			Usage:    "The influx server URL",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "org",
			Aliases:  []string{"o"},
			Usage:    "The influx org that owns the bucket",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "bucket",
			Aliases:  []string{"b"},
			Usage:    "The influx bucket for the measurements",
			Required: true,
		},
		&cli.IntFlag{
			Name:        "interval",
			Aliases:     []string{"i"},
			Usage:       "The interval in seconds between each probe",
			DefaultText: "10",
		},
	}

	app.Action = func(c *cli.Context) error {
		token := viper.GetString("INFLUX_TOKEN")
		if token == "" {
			return cli.Exit(errors.New("Environment variable INFLUX_TOKEN not set"), 1)
		}

		probe, err := NewProbe(c.String("server"), c.String("org"), c.String("bucket"), token)
		if err != nil {
			return cli.Exit(err, 1)
		}

		// Continuous loop to run the probe indefinitely
		for {
			err = probe.Run(c.String("destination"))
			if err != nil {
				log.Error().Err(err).Msg("Probing error")
			}
			// Wait for a short duration before probing again
			time.Sleep(time.Duration(c.Int("interval")) * time.Second) // Adjust the sleep duration as necessary
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(0)
	}
}
