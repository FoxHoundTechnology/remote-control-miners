package ant_miner_cgi

import (
	miner "foxhound/internal/application/miner/domain"
)

// cgi-bin/stats.cgi
// TODO: default values for response object

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
	Status    string      `json:"status"`
	Elapsed   int64       `json:"uptime"` // elapsed -> uptime
	Rate5s    float64     `json:"rate_5s"`
	Rate30m   float64     `json:"rate_30m"`
	RateAvg   float64     `json:"rate_avg"`
	RateIdeal float64     `json:"rate_ideal"`
	RateUnit  string      `json:"rate_unit"`
	MinerMode miner.Mode  `json:"miner_mode"`
	HashBoard []HashBoard `json:"hash_board"`
	Fan       []int       `json:"fan"`
	ChainNum  int         `json:"chain_num"` // for missing hash_board
}

type HashBoard struct {
	Index     int     `json:"index"`
	FreqAvg   int     `json:"freq_avg"`
	RateIdeal float64 `json:"rate_ideal"`
	TempPic   []int   `json:"temp_pic"`
	TempPcb   []int   `json:"temp_pcb"`
	TempChip  []int   `json:"temp_chip"`
}

func AntMinerCGIGetStats(ipAddress string) (*StatsResponse, error) {
	return &StatsResponse{}, nil
}
