package yt

import (
	utils "YTDownloaderCli/internal/common"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() {

	model := initialModel()

	// F, err := tea.LogToFile("debug.log", "debug")

	// utils.UtilError(err)

	// defer F.Close()

	p := tea.NewProgram(model, tea.WithMouseAllMotion())

	_, e := p.Run()
	utils.UtilError(e)

}
