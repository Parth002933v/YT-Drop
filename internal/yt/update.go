package yt

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	middlewares := []Middleware{handleQuitKeyMiddleware}
	runMiddlewareChain(middlewares, msg, &cmd, m)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch m.stage {
		case 1:
			USelectionStage(m, msg, &cmd)
		case 2:
			UUrlInputStage(m, msg, &cmd)
		case 3:
			UResposeDataStage(m, msg, &cmd)
		}

	default:
		var cmd tea.Cmd
		m.bubbles.spinner, cmd = m.bubbles.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}
