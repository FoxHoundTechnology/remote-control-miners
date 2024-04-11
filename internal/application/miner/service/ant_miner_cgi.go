package service

import (
	miner "foxhound/internal/application/miner/domain"
	queries "foxhound/internal/application/miner/queries/ant_miner_cgi"
)

type AntminerCGI struct {
	Miner       *miner.Miner
	Mode        *miner.Mode
	Status      *miner.Status
	Config      *miner.Config
	Stats       *miner.Stats
	Pools       *miner.Pool
	Temperature *miner.Temperature
	Fan         *miner.Fan
	// TODO: cgi-specific fields
}

// regular miner functions
func (a *AntminerCGI) SetNormalMode() error {
	return nil
}

func (a *AntminerCGI) SetSleepMode() error {
	return nil
}

func (a *AntminerCGI) SetLowPowerMode() error {
	return nil
}

func (a *AntminerCGI) CheckStatus() error {
	return nil
}

func (a *AntminerCGI) CheckConfig() error {
	return nil
}

func (a *AntminerCGI) CheckNetworkInfo() error {
	networkInfo, err := queries.AntMinerCGIGetNetworkInfo(a.Miner.IPAddress)
	if err != nil {
		return err
	}

	a.Miner.MacAddress = networkInfo.MacAddress
	a.Miner.IPAddress = networkInfo.IPAddress

	return nil
}

func (a *AntminerCGI) CheckSystemInfo() error {
	return nil
}

func (a *AntminerCGI) ChangePool() error {
	return nil
}

func (a *AntminerCGI) Reboot() error {
	return nil
}
