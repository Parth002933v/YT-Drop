package yt

import (
	"YTDownloaderCli/internal/utils"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func USelectionStage(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case tea.KeyCtrlC.String():
			return m, tea.Quit

		case tea.KeyUp.String():
			if m.contentTypeSelection.cursor > 0 {
				m.contentTypeSelection.cursor--
			}

		case tea.KeyDown.String():
			if m.contentTypeSelection.cursor < len(m.contentTypeSelection.choices)-1 {
				m.contentTypeSelection.cursor++
			}

		case tea.KeyEnter.String():

			if m.contentTypeSelection.cursor > 0 {
				m.isPlaylist = true
			}

			m.stage++
		}
	}
	return m, nil
}

func UUrlInputStage(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	current := &m.questions[m.index]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case tea.KeyCtrlC.String():
			return m, tea.Quit

		case tea.KeyTab.String():
			current.input.SetValue(current.defaultVal)
			return m, nil

		case tea.KeyEnter.String():

			if current.input.Value() != "" {
				m.isLoading = true

				switch m.contentTypeSelection.cursor {
				case 0:
					go func() {
						detail := m.client.GetVideoDetail(current.input.Value())
						m.isLoading = false
						m.stage++
						utils.SetRequiredVideoFormat(detail)
						utils.SetQualitySelectionChoiceValue(&m.QualitySelection.choices, detail.Formats)
						m.data = VideoResData{video: *detail}
					}()
				case 1:
					go func() {
						detail := m.client.GetVideoPlaylistDetail(current.input.Value())
						m.isLoading = false
						m.stage++
						m.data = VideoResData{playlist: *detail}

					}()
				}
				current.input.Blur()
			}
			return m, cmd
		}

	default:
		var cmd tea.Cmd
		m.bubbles.spinner, cmd = m.bubbles.spinner.Update(msg)
		return m, cmd
	}

	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func UResposeDataStage(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case tea.KeyCtrlC.String():
			return m, tea.Quit

		case tea.KeyUp.String():
			if m.QualitySelection.cursor > 0 {
				m.QualitySelection.cursor--
			}

		case tea.KeyDown.String():
			if m.QualitySelection.cursor < len(m.QualitySelection.choices)-1 {
				m.QualitySelection.cursor++
			}

		case tea.KeyEnter.String():
			m.stage++

			go func() {
				//stream :  io.ReadCloser  //size : int64 //err: error
				stream, s, err := m.client.GetDownloadStream(&m.data.video, &m.data.video.Formats[m.QualitySelection.cursor])
				utils.PError(err)

				defer stream.Close()

				file, err := os.Create(fmt.Sprintf("%s.mp4", utils.SanitizeFileName(m.data.video.Title)))
				utils.PError(err)

				defer file.Close()

				pw := &ProgressWriter{
					Writer:     file,
					TotalBytes: s,
					ProgressDisplay: func(progress float64) {
						m.data.downloadPrecentage = progress
					},
				}

				_, err = io.Copy(pw, stream)
				utils.PError(err)
			}()

		}
	}
	return m, cmd
}
