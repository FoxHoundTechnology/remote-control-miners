package dashboard

import (
	"time"

	"github.com/FoxHoundTechnology/remote-control-miners/terminal/query"
	"github.com/charmbracelet/bubbles/table"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	queryClient  query.Query
	tickDuration time.Duration

	miners  []query.MinerData
	table   table.Model
	sidebar string
}

func InitModel(query query.Query, updateInterval time.Duration) (*Model, error) {
	// Interval query won't start immediately when the client runs.
	// This leaves MinerData empty until client clock syncs with system clock.
	// Bypass this issue by getting the data when client starts.
	// https://github.com/charmbracelet/bubbletea/blob/v0.26.1/commands.go#L96

	resp, err := query.GetMiners()
	if err != nil {
		return nil, err
	}

	table := newMinerTable()
	table.SetStyles(minerTableStyle())
	sidebar := ""
	// sidebar.SetStyle

	model := Model{
		queryClient:  query,
		tickDuration: updateInterval,
		miners:       resp.Data,

		table:   table,
		sidebar: sidebar,
	}

	model.updateMinerTable()
	model.updateMinerStatsSidebar()

	return &model, nil
}

func (m Model) Init() tea.Cmd {
	return m.tickQueryMiner()
}
