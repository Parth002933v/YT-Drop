package yt

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case tea.KeyCtrlC.String():
			return m, tea.Quit

		default:

			switch m.stage {
			case 1:
				USelectionStage(m, msg)

			case 2:
				UUrlInputStage(m, msg)
			case 3:
				UResposeDataStage(m, msg)

			default:
				m.stage = 0
				USelectionStage(m, msg)
			}
		}

	default:
		var cmd tea.Cmd
		m.bubbles.spinner, cmd = m.bubbles.spinner.Update(msg)
		return m, cmd

	}

	return m, cmd
}
