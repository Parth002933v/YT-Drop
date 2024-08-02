package ffmpeg

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Middleware func(tea.Msg, *tea.Cmd, *ffmpegModel) bool

func HandleQuitKeyMiddleware(msg tea.Msg, cmd *tea.Cmd, m *ffmpegModel) bool {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == tea.KeyCtrlC.String() {
		m.questions.input.Blur()
		return false
	}
	return true
}

func RunMiddlewareChain(middlewares []Middleware, msg tea.Msg, cmd *tea.Cmd, m *ffmpegModel) {
	for _, middleware := range middlewares {
		if !middleware(msg, cmd, m) {
			return
		}
	}
}
