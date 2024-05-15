package service

import (
	"fmt"
	"sync"

	commands "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/commands"
	queries "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/ant_miner_cgi/queries"
	domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
)

// TODO: add the logic for updating the pool "stats"
// TODO: custom errors & logger for miner service
// TODO: add contexts of miner API responses
// TODO: add comments
// TODO: timeout logic for set_miner_conf requests

type AntminerCGI struct {
	Miner       domain.Miner
	Mode        domain.Mode
	Status      domain.Status
	Config      domain.Config
	Stats       domain.Stats
	Pools       []domain.Pool
	Temperature []domain.TemperatureSensor
	Fan         []domain.FanSensor
	FanCtrl     bool   // fan control enabled/disabled
	FanPwm      string // fan pwm value
	FreqLevel   string // frequency level
	Model       string // miner model name (e.g. S19, S17)
	rwMutex     *sync.RWMutex
}

func NewAntminerCGI(config domain.Config, miner domain.Miner, modelName string) *AntminerCGI {
	return &AntminerCGI{
		Miner:  miner,
		Mode:   domain.SleepMode,
		Status: domain.Online,
		Config: config,
		Stats: domain.Stats{
			HashRate:  0,
			RateIdeal: 0,
			Uptime:    0,
		},
		Pools:       make([]domain.Pool, 0),
		Temperature: make([]domain.TemperatureSensor, 0),
		Fan:         make([]domain.FanSensor, 0),
		FanCtrl:     true,
		FanPwm:      "100",
		FreqLevel:   "",
		Model:       modelName,
		rwMutex:     new(sync.RWMutex),
	}
}

func (a *AntminerCGI) CheckConfig() error {
	GetMinerConfigResponse, err := queries.AntMinerCGIGetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.FanCtrl = GetMinerConfigResponse.BitmainFanCtrl
	a.FanPwm = GetMinerConfigResponse.BitmainFanPWM
	a.FreqLevel = GetMinerConfigResponse.FreqLevel

	// TODO! R&D
	// a.Pools = GetMinerConfigResponse.Pools ???

	if GetMinerConfigResponse.MinerMode == 3 {
		a.Mode = domain.Mode(domain.LowPowerMode)
	} else {
		a.Mode = domain.Mode(GetMinerConfigResponse.MinerMode)
	}

	fmt.Println("new miner mode in service", domain.Mode(GetMinerConfigResponse.MinerMode))
	return nil
}

func (a *AntminerCGI) SetNormalMode() error {
	a.CheckConfig()

	SetMinerConfigResponse, err := commands.AntminerCGISetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress, commands.SetMinerConfigPayload{
		BitmainFanCtrl: a.FanCtrl,
		BitmainFanPWM:  a.FanPwm,
		FreqLevel:      a.FreqLevel,
		MinerMode:      "0", // Normal Mode
		Pools:          a.Pools,
	})
	if err != nil {
		return err
	}
	if SetMinerConfigResponse.Stats != "success" {
		return fmt.Errorf("failed to set miner config")
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.Mode = domain.NormalMode

	return nil
}

func (a *AntminerCGI) SetSleepMode() error {
	a.CheckConfig()

	SetMinerConfigResponse, err := commands.AntminerCGISetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress, commands.SetMinerConfigPayload{
		BitmainFanCtrl: a.FanCtrl,
		BitmainFanPWM:  a.FanPwm,
		FreqLevel:      a.FreqLevel,
		MinerMode:      "1", // Sleep Mode
		Pools:          a.Pools,
	})
	if err != nil {
		return err
	}

	if SetMinerConfigResponse.Stats != "success" {
		return fmt.Errorf("failed to set miner config")
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.Mode = domain.SleepMode

	return nil
}

func (a *AntminerCGI) SetLowPowerMode() error {
	a.CheckConfig()

	SetMinerConfigResponse, err := commands.AntminerCGISetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress, commands.SetMinerConfigPayload{
		BitmainFanCtrl: a.FanCtrl,
		BitmainFanPWM:  a.FanPwm,
		FreqLevel:      a.FreqLevel,
		MinerMode:      "3", // Low Power Mode
		Pools:          a.Pools,
	})
	if err != nil {
		return err
	}

	if SetMinerConfigResponse.Stats != "success" {
		return fmt.Errorf("failed to set miner config")
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.Mode = domain.LowPowerMode

	return nil
}

func (a *AntminerCGI) CheckStats() error {
	GetStatsResponse, err := queries.AntMinerCGIGetStats(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.Stats = domain.Stats{
		HashRate:  GetStatsResponse.Rate5s,
		RateIdeal: GetStatsResponse.RateIdeal,
		Uptime:    int(GetStatsResponse.Elapsed),
	}
	a.Mode = domain.Mode(GetStatsResponse.Mode)

	for index, tempSensor := range GetStatsResponse.Chain {
		pcbSensors := []int{}
		pcbSensors = append(pcbSensors, tempSensor.TempPcb...)
		a.Temperature = append(a.Temperature, domain.TemperatureSensor{
			Name:       fmt.Sprintf("Chain %d", index),
			PcbSensors: pcbSensors,
		})
	}

	for index, speed := range GetStatsResponse.Fan {
		a.Fan = append(a.Fan, domain.FanSensor{
			Name:  fmt.Sprintf("Fan %d", index),
			Speed: speed,
		})
	}

	return nil
}

func (a *AntminerCGI) CheckPools() error {
	GetPoolsResponse, err := queries.AntMinerCGIGetPools(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	fmt.Println("pool elements before for loop in CheckPools :", GetPoolsResponse)

	a.Pools = make([]domain.Pool, len(GetPoolsResponse))

	for index, pool := range GetPoolsResponse {
		a.Pools[index] = domain.Pool{
			Status:   pool.Status,
			Accepted: pool.Accepted,
			Rejected: pool.Rejected,
			Stale:    pool.Stale,
		}
	}

	fmt.Println("a.Pools element after for loop in checkpools", a.Pools)

	return nil
}

func (a *AntminerCGI) CheckSystemInfo() error {
	GetSystemInfoResponse, err := queries.AntMinerCGIGetSystemInfo(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.Config.Firmware = GetSystemInfoResponse.FirmwareType
	a.Miner.IPAddress = GetSystemInfoResponse.IPAddress
	a.Miner.MacAddress = GetSystemInfoResponse.MacAddress
	a.Model = GetSystemInfoResponse.MinerType // NOTE: miner model name

	return nil
}

func (a *AntminerCGI) Reboot() error {
	err := queries.AntMinerCGIReboot(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	return nil
}

func (a *AntminerCGI) ChangePool(pools []domain.Pool) error {
	a.CheckConfig()
	SetMinerConfigResponse, err := commands.AntminerCGISetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress, commands.SetMinerConfigPayload{
		BitmainFanCtrl: a.FanCtrl,
		BitmainFanPWM:  a.FanPwm,
		FreqLevel:      a.FreqLevel,
		MinerMode:      fmt.Sprintf("%d", a.Mode), // Normal Mode
		Pools:          pools,
	})
	if err != nil {
		return err
	}

	if SetMinerConfigResponse.Stats != "success" {
		return fmt.Errorf("failed to set miner config")
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.Pools = pools

	return nil
}

func (a *AntminerCGI) CheckNetworkInfo() error {
	GetNetWorkInfoResponse, err := queries.AntMinerCGIGetNetworkInfo(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	a.Miner.IPAddress = GetNetWorkInfoResponse.IPAddress
	a.Miner.MacAddress = GetNetWorkInfoResponse.MacAddress

	return nil
}
