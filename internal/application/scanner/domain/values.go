package domain

type MinerType int

const (
	AntminerCgi MinerType = iota // 0
	//...
)

type Config struct {
	Interval int // in minutes
	Username string
	Password string
}

// threshold value for alert (i.e. 100 TH/s, 70 degrees Celsius, etc.)
type AlertTriggerValue int

// threshold count (i.e. count as a number of machine, rate)
type AlertMachineCount int

// threshold type (i.e. "# of machines: 0 ~ X" or "% of machines: 0 ~ 100" in a fleet)
// depending on the threshold type below
type AlertThresholdType int

const (
	Count AlertThresholdType = iota // 0
	Rate                            // 1
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
	Reboot AlertActionType = iota // 0
	Sleep                         // 1
	Normal                        // 2
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
