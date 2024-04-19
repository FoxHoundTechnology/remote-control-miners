package domain

type MinerType string

const (
	AntminerCgi MinerType = "antminer_cgi"
	//...
)

type Config struct {
	Interval int // in minutes
	Username string
	Password string
}

type AlertValue int

type AlertThresholdType int

const (
	ThresholdCount AlertThresholdType = iota // 0
	ThresholdRate                            // 1
)

type AlertConditionType int

const (
	Hashrate          AlertConditionType = iota // 0
	Temperature                                 // 1
	FanSpeed                                    // 2
	PoolShares                                  // 3
	OfflineMiners                               // 4
	MissingHashboards                           // 5
)

type AlertActionType int

const (
	Reboot     AlertActionType = iota // 0
	Sleep                             // 1
	Normal                            // 2
	ChangePool                        // 3
)

type AlertLayerType int

const (
	InfoAlert    AlertLayerType = iota // 0
	WarningAlert                       // 1
	ErrorAlert                         // 2
	FatalAlert                         // 3
)

type AlertState int

const (
	Monitoring AlertState = iota // 0
	Triggered                    // 1
	Resolving                    // 2
	Resolved                     // 3
)

type Log string
