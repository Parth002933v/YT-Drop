package ffmpeg

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *ffmpegModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	middlewares := []Middleware{HandleQuitKeyMiddleware}
	RunMiddlewareChain(middlewares, msg, &cmd, m)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch m.stage {

		case 1:
			USelectionStage(m, msg, &cmd)

		case 2:
			UPostSelection(m, msg, &cmd)

		}

	default:
		return m, cmd
	}

	return m, cmd

}
