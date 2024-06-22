package dashboard

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type queryMinerErrMsg error
type tickQueryMinerMsg time.Time

func (m *Model) tickQueryMiner() tea.Cmd {
	return tea.Every(m.tickDuration, func(t time.Time) tea.Msg {
		resp, err := m.queryClient.GetMiners()
		if err != nil {
			return queryMinerErrMsg(err)
		}
		m.miners = resp.Data
		return tickQueryMinerMsg(t)
	})
}
