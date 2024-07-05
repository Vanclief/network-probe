package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vanclief/ez"
	"github.com/vanclief/network-probe/mtr"
)

type Probe struct{}

func NewProbe() (*Probe, error) {
	const op = "NewProbe"

	p := &Probe{}

	return p, nil
}

func (p *Probe) Run(destination string) error {
	const op = "Probe.Run"

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// TODO: Go back to 10
	cmd := exec.Command("mtr", "-r", "-c", "1", "--json", destination)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return ez.Wrap(op, err)
	}

	log.Info().Str("destination", destination).Msg("Ran MTR")

	var res mtr.OutputJSON

	err = json.Unmarshal(stdout.Bytes(), &res)
	if err != nil {
		return ez.Wrap(op, err)
	}

	res.Report.Timestamp = time.Now().Unix()

	fmt.Println("Res", res.Report.String())

	return nil
}
