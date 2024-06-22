package dashboard

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	// Render the table and sidebar individually with their styles
	tableRendered := tableStyle.Render(m.table.View())
	sidebarRendered := sidebarStyle.Render(m.sidebar)

	// Use lipgloss JoinHorizontal to properly place blocks side by side
	// Adjusting to 'lipgloss.Top' aligns both components at the top
	fullView := lipgloss.JoinHorizontal(lipgloss.Top, tableRendered, "  ", sidebarRendered)

	// Return the combined view
	return fullView + "\n"
}
