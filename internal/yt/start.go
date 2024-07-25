package yt

import (
	"YTDownloaderCli/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() {

	model := initialModel()

	f, err := tea.LogToFile("debug.log", "debug")

	utils.PError(err)

	defer f.Close()

	p := tea.NewProgram(model, tea.WithMouseAllMotion())

	_, e := p.Run()
	utils.PError(e)

}
