package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rs/zerolog/log"
	"github.com/vanclief/ez"
	"github.com/vanclief/network-probe/mtr"
)

type Probe struct {
	InfluxClient influxdb2.Client
	InfluxOrg    string
	InfluxBucket string
}

func NewProbe(influxURL, org, bucket, token string) (*Probe, error) {
	const op = "NewProbe"

	p := &Probe{
		InfluxOrg:    org,
		InfluxBucket: bucket,
		InfluxClient: influxdb2.NewClient(influxURL, token),
	}

	log.Info().
		Str("bucket", bucket).
		Str("org", org).
		Str("url", influxURL).
		Msg("Created probe")

	return p, nil
}

func (p *Probe) Run(destination, protocol string) error {
	const op = "Probe.Run"

	if protocol != "IPv4" && protocol != "IPv6" {
		return ez.New(op, ez.EINVALID, "Invalid protocol", nil)
	}

	report, err := p.GenerateReport(destination, protocol)
	if err != nil {
		return ez.Wrap(op, err)
	}

	err = p.Send(report, protocol)
	if err != nil {
		return ez.Wrap(op, err)
	}

	return nil
}

func (p *Probe) GenerateReport(destination string, protocol string) (*mtr.Report, error) {
	const op = "Probe.GenerateReport"

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var protocolArg string

	switch protocol {
	case "IPv4":
		protocolArg = "-4"
	case "IPv6":
		protocolArg = "-6"
	default:
		return nil, ez.New(op, ez.EINVALID, "Invalid protocol", nil)
	}

	cmd := exec.Command("mtr", protocolArg, "-r", "-c", "10", "--json", destination)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	log.Info().
		Str("destination", destination).
		Str("protocol", protocol).
		Msg("Ran MTR")

	var res mtr.OutputJSON

	err = json.Unmarshal(stdout.Bytes(), &res)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	res.Report.Timestamp = time.Now().Unix()

	return &res.Report, nil
}

func (p *Probe) Send(report *mtr.Report, protocol string) error {
	const op = "Probe.Send"

	// Get non-blocking write client
	writeAPI := p.InfluxClient.WriteAPIBlocking(p.InfluxOrg, p.InfluxBucket)

	// Create a new point using full params constructor
	for i, hub := range report.Hubs {

		point := influxdb2.NewPoint(protocol,
			map[string]string{"source": report.MTR.Source, "destination": report.MTR.Destination, "host": hub.Host, "hop": strconv.Itoa(i)},
			map[string]interface{}{
				"number_of_tests": report.MTR.NumberOfTests,
				"pattern_size":    report.MTR.PatternSize,
				"bit_pattern":     report.MTR.BitPattern,
				"count":           hub.Count,
				"loss":            hub.LossPercent,
				"sent":            hub.Sent,
				"last":            hub.Last,
				"avg":             hub.Average,
				"best":            hub.Best,
				"worst":           hub.Worst,
				"std":             hub.StandardDeviation,
			},
			time.Unix(report.Timestamp, 0))

		err := writeAPI.WritePoint(context.Background(), point)
		if err != nil {
			log.Error().Err(err).Msg("Error writing point")
		}
	}

	// Write the point
	err := writeAPI.Flush(context.Background())
	if err != nil {
		return ez.Wrap(op, err)
	}

	log.Info().
		Str("destination", report.MTR.Destination).
		Str("protocol", protocol).
		Str("source", report.MTR.Source).
		Msg("Written to InfluxDB")

	return nil
}
