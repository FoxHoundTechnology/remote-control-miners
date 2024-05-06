package domain

// TODO: validation
// TODO: create custom response objects for domain entities
type Miner struct {
	MacAddress string
	IPAddress  string
}

type MinerController interface {
	SetNormalMode() error
	SetSleepMode() error
	SetLowPowerMode() error

	CheckStats() error
	CheckPools() error
	CheckConfig() error

	CheckNetworkInfo() error
	CheckSystemInfo() error

	ChangePool([]Pool) error
	Reboot() error
}
