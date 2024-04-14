package service

import (
	"fmt"
	"time"

	commands "foxhound/internal/application/miner/ant_miner_cgi/commands"
	queries "foxhound/internal/application/miner/ant_miner_cgi/queries"
	domain "foxhound/internal/application/miner/domain"
)

// TODO: custom errors for miner service
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
}

func (a *AntminerCGI) CheckConfig() error {
	GetMinerConfigResponse, err := queries.AntMinerCGIGetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.FanCtrl = GetMinerConfigResponse.BitmainFanCtrl
	a.FanPwm = GetMinerConfigResponse.BitmainFanPWM
	a.FreqLevel = GetMinerConfigResponse.FreqLevel
	a.Pools = GetMinerConfigResponse.Pools
	a.Mode = domain.Mode(GetMinerConfigResponse.MinerMode)

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

	return nil
}

func (a *AntminerCGI) SetSleepMode() error {
	a.CheckConfig()

	SetMinerConfigResponse, err := commands.AntminerCGISetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress, commands.SetMinerConfigPayload{
		BitmainFanCtrl: a.FanCtrl,
		BitmainFanPWM:  a.FanPwm,
		FreqLevel:      a.FreqLevel,
		MinerMode:      "1", // Normal Mode
		Pools:          a.Pools,
	})
	if err != nil {
		return err
	}

	if SetMinerConfigResponse.Stats != "success" {
		return fmt.Errorf("failed to set miner config")
	}

	return nil
}

func (a *AntminerCGI) SetLowPowerMode() error {
	a.CheckConfig()

	SetMinerConfigResponse, err := commands.AntminerCGISetMinerConfig(a.Config.Username, a.Config.Password, a.Miner.IPAddress, commands.SetMinerConfigPayload{
		BitmainFanCtrl: a.FanCtrl,
		BitmainFanPWM:  a.FanPwm,
		FreqLevel:      a.FreqLevel,
		MinerMode:      "3", // Sleep Mode
		Pools:          a.Pools,
	})
	if err != nil {
		return err
	}

	if SetMinerConfigResponse.Stats != "success" {
		return fmt.Errorf("failed to set miner config")
	}

	return nil
}

func (a *AntminerCGI) CheckStatus() error {
	return nil
}

func (a *AntminerCGI) CheckStats() error {
	GetStatsResponse, err := queries.AntMinerCGIGetStats(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.Stats = domain.Stats{
		HashRate:    GetStatsResponse.Rate5s,
		RateIdeal:   GetStatsResponse.RateIdeal,
		Uptime:      int(GetStatsResponse.Elapsed),
		LastUpdated: time.Now(),
	}
	a.Mode = domain.Mode(GetStatsResponse.Mode)

	for index, tempereture := range GetStatsResponse.Chain {
		a.Temperature = append(a.Temperature, domain.TemperatureSensor{
			Name:    fmt.Sprintf("Chain %d", index),
			TempPcb: tempereture.TempPcb,
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

func (a *AntminerCGI) CheckSystemInfo() error {
	GetSystemInfoResponse, err := queries.AntMinerCGIGetSystemInfo(a.Config.Username, a.Config.Password, a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.Config.Firmware = GetSystemInfoResponse.FirmwareType
	a.Miner.IPAddress = GetSystemInfoResponse.IPAddress
	a.Miner.MacAddress = GetSystemInfoResponse.MacAddress

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

	return nil
}

func (a *AntminerCGI) CheckNetworkInfo() error {
	return nil
}
