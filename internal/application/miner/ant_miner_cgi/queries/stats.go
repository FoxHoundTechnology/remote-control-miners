package queries

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	miner "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"

	"github.com/sirupsen/logrus"
)

// cgi-bin/stats.cgi
// TODO: default values for response object
// TODO: missing hash board

type rawGetStatsResponse struct {
	Status Status  `json:"STATUS"`
	Stats  []Stats `json:"STATS"`
}

type Stats struct {
	Elapsed   int64   `json:"elapsed"`
	Rate5s    float64 `json:"rate_5s"`
	Rate30m   float64 `json:"rate_30m"`
	RateAvg   float64 `json:"rate_avg"`
	RateIdeal float64 `json:"rate_ideal"`
	RateUnit  string  `json:"rate_unit"`
	ChainNum  int     `json:"chain_num"`
	FanNum    int     `json:"fan_num"`
	Fan       []int   `json:"fan"`
	HwpTotal  float64 `json:"hwp_total"`
	MinerMode int     `json:"miner-mode"`
	FreqLevel int     `json:"freq-level"`
	Chain     []Chain `json:"chain"`
}

type Status struct {
	Status string `json:"STATUS"`
}

type Chain struct {
	Index        int     `json:"index"`
	FreqAvg      int     `json:"freq_avg"`
	RateIdeal    float64 `json:"rate_ideal"`
	RateReal     float64 `json:"rate_real"`
	AsicNum      int     `json:"asic_num"`
	Asic         string  `json:"asic"`
	TempPic      []int   `json:"temp_pic"`
	TempPcb      []int   `json:"temp_pcb"`
	TempChip     []int   `json:"temp_chip"`
	Hw           int     `json:"hw"`
	EepromLoaded bool    `json:"eeprom_loaded"`
	Sn           string  `json:"sn"`
	Hwp          float64 `json:"hwp"`
}

type StatsResponse struct {
	Status  string  `json:"status"`
	Elapsed int64   `json:"uptime"` // elapsed -> uptime
	Rate5s  float64 `json:"rate_5s"`
	// Rate30m   float64     `json:"rate_30m"`
	// RateAvg   float64     `json:"rate_avg"`
	RateIdeal float64    `json:"rate_ideal"`
	RateUnit  string     `json:"rate_unit"`
	Mode      miner.Mode `json:"miner_mode"`
	Chain     []Chain    `json:"chain"`
	Fan       []int      `json:"fan"`
	ChainNum  int        `json:"chain_num"` // for missing hash_board
}

func AntMinerCGIGetStats(username, password, ipAddress string) (*StatsResponse, error) {

	t := http_auth.NewTransport(username, password)
	newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/stats.cgi", ipAddress), nil)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":      err,
			"newRequest": newRequest,
		}).Info("Error creating new request")

		return nil, err
	}

	resp, err := t.RoundTrip(newRequest)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err,
			"resp":  resp,
		}).Info("Error creating new request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"body":  body,
		}).Info("Error reading response body")
		return nil, err
	}

	var rawGetStatsResponse rawGetStatsResponse
	err = json.Unmarshal(body, &rawGetStatsResponse)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"reponse": body,
		}).Info("Error unmarshalling response body")
		return nil, err
	}

	if len(rawGetStatsResponse.Stats) == 0 {
		fmt.Println("EDGE CASE ---------------------")
		fmt.Println("ip address", ipAddress)

		return &StatsResponse{
			Status:    "S",
			Elapsed:   0,
			RateIdeal: 0,
			Rate5s:    0,
			RateUnit:  "GH/s",
			Mode:      miner.SleepMode,
			Chain: []Chain{
				{
					Index:   0,
					TempPcb: []int{0},
					// NOTE: fallback values for additional fields go here
				},
			},
			Fan: []int{
				0,
			},
			ChainNum: 0,
		}, nil
	}

	mode := 0
	switch rawGetStatsResponse.Stats[0].MinerMode {
	case 0:
		mode = 0
	case 1:
		mode = 1
	case 3: // low power mode
		mode = 2
	}

	return &StatsResponse{
		Status:    rawGetStatsResponse.Status.Status,
		Elapsed:   rawGetStatsResponse.Stats[0].Elapsed,
		RateIdeal: rawGetStatsResponse.Stats[0].RateIdeal,
		Rate5s:    rawGetStatsResponse.Stats[0].Rate5s,
		RateUnit:  rawGetStatsResponse.Stats[0].RateUnit,
		Mode:      miner.Mode(mode),
		Chain:     rawGetStatsResponse.Stats[0].Chain,
		Fan:       rawGetStatsResponse.Stats[0].Fan,
		ChainNum:  rawGetStatsResponse.Stats[0].ChainNum,
	}, nil
}
