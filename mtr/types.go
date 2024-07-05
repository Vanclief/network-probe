package mtr

import (
	"fmt"
	"strings"
	"time"
)

type OutputJSON struct {
	Report Report `json:"report"`
}

type Report struct {
	MTR       MTR   `json:"mtr"`
	Hubs      []HUB `json:"hubs"`
	Timestamp int64 `json:"timestamp"`
}

type MTR struct {
	Source        string `json:"src"`
	Destination   string `json:"dst"`
	TypeOfService int    `json:"tos"`
	NumberOfTests int    `json:"tests"`
	PatternSize   string `json:"psize"`
	BitPattern    string `json:"bitpattern"`
}

type HUB struct {
	Count             int     `json:"count"`
	Host              string  `json:"host"`
	LossPercent       float32 `json:"Loss%"`
	Sent              int     `json:"Snt"`
	Last              float32 `json:"Last"`
	Average           float32 `json:"Avg"`
	Best              float32 `json:"Best"`
	Worst             float32 `json:"Wrst"`
	StandardDeviation float32 `json:"StDev"`
}

func (r *Report) String() string {
	var sb strings.Builder
	sb.WriteString("Report Timestamp: ")
	sb.WriteString(time.Unix(r.Timestamp, 0).String())
	sb.WriteString("\nMTR Details:\n")
	sb.WriteString(fmt.Sprintf("Source: %s, Destination: %s, Type of Service: %d, Number of Tests: %s, Pattern Size: %s, Bit Pattern: %s\n",
		r.MTR.Source, r.MTR.Destination, r.MTR.TypeOfService, r.MTR.NumberOfTests, r.MTR.PatternSize, r.MTR.BitPattern))
	sb.WriteString("Hubs Details:\n")
	for _, hub := range r.Hubs {
		sb.WriteString(fmt.Sprintf("Host: %s, Count: %d, Loss%%: %.2f, Sent: %d, Last: %.2f, Avg: %.2f, Best: %.2f, Worst: %.2f, StDev: %.2f\n",
			hub.Host, hub.Count, hub.LossPercent, hub.Sent, hub.Last, hub.Average, hub.Best, hub.Worst, hub.StandardDeviation))
	}
	return sb.String()
}
