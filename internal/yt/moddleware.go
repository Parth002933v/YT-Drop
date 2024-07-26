package yt

import tea "github.com/charmbracelet/bubbletea"

type Middleware func(tea.Msg, *tea.Cmd, *model) bool

func handleQuitKeyMiddleware(msg tea.Msg, cmd *tea.Cmd, m *model) bool {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == tea.KeyCtrlC.String() {
		*cmd = tea.Quit
		return false // Stop further processing
	}
	return true // Continue processing
}

func runMiddlewareChain(middlewares []Middleware, msg tea.Msg, cmd *tea.Cmd, m *model) {
	for _, middleware := range middlewares {
		if !middleware(msg, cmd, m) {
			return
		}
	}
}
