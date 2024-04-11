package miner

// TODO: validation
type Miner struct {
	MacAddress string
	IPAddress  string
	Owner      string
}

type MinerController interface {
	SetNormalMode() error
	SetSleepMode() error
	SetLowPowerMode() error

	CheckStatus() error
	CheckConfig() error
	CheckNetworkInfo() error
	CheckSystemInfo() error

	ChangePool() error
	Reboot() error
}
