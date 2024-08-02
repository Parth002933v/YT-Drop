package ffmpeg

import (
	utils "YTDownloaderCli/internal/common"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() {

	model := initialModel()

	f, err := tea.LogToFile("debug.log", "debug")

	utils.UtilError(err)

	defer f.Close()

	p := tea.NewProgram(model)

	_, e := p.Run()
	utils.UtilError(e)

}


