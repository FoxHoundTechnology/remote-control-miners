package dashboard

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"

	"github.com/FoxHoundTechnology/remote-control-miners/terminal/query"
)

func newMinerTable() table.Model {
	columns := []table.Column{
		{Title: "MAC Addr", Width: 18},
		{Title: "Miner Model", Width: 16},
		{Title: "Hashrate", Width: 10},
		{Title: "Status", Width: 6},
		{Title: "Mode", Width: 11},
		{Title: "Uptime", Width: 14},
		{Title: "Fleet ID", Width: 10},
		// NOTE: hidden detail info for each miner can be set with empty string
		{Title: "", Width: 0}, // miner config, pool, etc
		// ...
	}

	rows := make([]table.Row, 0)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(30),
	)
	return t
}

func newMinerTableRowFromData(data []query.MinerData) []table.Row {
	// NOTE: append 2 as # of extra hidden fields
	rows := make([]table.Row, len(data)+2)

	for i, miner := range data {
		rows[i] = []string{
			miner.Miner.MacAddress,
			miner.ModelName,
			hashRateInTHToString(miner.Stats.HashRate),
			minerStatusToString(miner.Status),
			minerModeToString(miner.Mode),
			minerUpTimeToString(miner.Stats.Uptime),
			minerFleetIDToString(miner.FleetID),
		}
	}

	return rows
}

func newMinerStatsSidebarFromData(data []query.MinerData) string {

	Hashrate := 0.0        // sum of actual hashrate for each model
	NameplateHashrate := 0 // sum of expected hashrate for each model
	OnlineCount := 0
	OfflineCount := 0
	TotalCount := len(data)

	for _, miner := range data {
		Hashrate += miner.Stats.HashRate
		NameplateHashrate += miner.Stats.RateIdeal
		if miner.Status == 0 {
			OnlineCount++
			// other status codes can go here
		} else {
			OfflineCount++
		}
	}

	newSidebar := fmt.Sprintf("Hashrate: %s PH/s\n"+
		"Nameplate Hashrate: %s PH/s\n"+
		"Online: %s\n"+
		"Offline: %s\n"+
		"Total: %s\n",
		hashRateInPHToString(Hashrate),
		nameplateHashRateInPHToString(NameplateHashrate),
		countToString(OnlineCount),
		countToString(OfflineCount),
		countToString(TotalCount),
	)

	return newSidebar
}

// TH/s
func hashRateInTHToString(rate float64) string {
	rateInTHs := rate / 1000.0

	return fmt.Sprintf("%.2f", rateInTHs)
}

// PH/s
func hashRateInPHToString(rate float64) string {
	rateInPHs := rate / 1000000.0

	return fmt.Sprintf("%.2f", rateInPHs)
}

func nameplateHashRateInPHToString(rate int) string {
	rateInPHs := rate / 1000000

	return fmt.Sprintf("%d", rateInPHs)
}

func countToString(count int) string {
	return fmt.Sprintf("%d", count)
}

func minerStatusToString(status int) string {
	return fmt.Sprintf(statusLabelMap[status])
}

func minerModeToString(mode int) string {
	return fmt.Sprintf(modeLabelMap[mode])
}

func minerUpTimeToString(uptime int) string {
	// TBD: where to put conversion logic
	duration := time.Duration(uptime) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02dh:%02dm:%02ds", hours, minutes, seconds)
}

func minerFleetIDToString(id int) string {
	return fmt.Sprintf("%d", id)
}
