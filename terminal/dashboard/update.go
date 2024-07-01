package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("%v selected!", m.table.SelectedRow()),
			)
		}
	case tickQueryMinerMsg:
		m.updateMinerTable()
		m.updateMinerStatsSidebar()

		return m, m.tickQueryMiner()
	case queryMinerErrMsg:
		// NOTE: do something about it
		return m, tea.Quit
	}

	//m.sidebar = "Upon my head they've placed a fruitless crown and a barren sceptre in my grip."
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// update table every query
func (m *Model) updateMinerTable() {
	m.table.SetRows(newMinerTableRowFromData(m.miners))
}

func (m *Model) updateMinerStatsSidebar() {
	m.sidebar = newMinerStatsSidebarFromData(m.miners)
}
