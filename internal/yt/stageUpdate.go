package yt

import (
	"YTDownloaderCli/internal/utils"
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kkdai/youtube/v2"
)

func USelectionStage(m *model, msg tea.Msg, cmd *tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

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
}

func UUrlInputStage(m *model, msg tea.Msg, cmd *tea.Cmd) {

	current := &m.questions[m.index]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case tea.KeyTab.String():
			current.input.SetValue(current.defaultVal)

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
		}

	default:
		m.bubbles.spinner, *cmd = m.bubbles.spinner.Update(msg)
	}

	current.input, *cmd = current.input.Update(msg)
}

func UResposeDataStage(m *model, msg tea.Msg, cmd *tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

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

			go downloadAndMerge(m)

		}
	}
}

func downloadAndMerge(m *model) {

	format := m.data.video.Formats[m.QualitySelection.cursor]
	if strings.Contains(format.MimeType, "avc1") {

		download(m, &format, m.data.video.Title, format.QualityLabel, "mp4")
		m.data.downloadPrecentage = float64(0)

		download(m, utils.GetMaxAudioQuality(m.data.video.Formats), m.data.video.Title, format.AudioQuality, "mp3")

	} else {
		download(m, utils.GetMaxAudioQuality(m.data.video.Formats), m.data.video.Title, format.AudioQuality, "mp3")
	}

}

func download(m *model, format *youtube.Format, fileName string, surfixName string, extention string) {
	//* start downlod
	stream, s, err := m.client.GetDownloadStream(&m.data.video, format)
	utils.PError(err)

	defer stream.Close()

	file, err := os.Create(fmt.Sprintf("%s_%s.%s", utils.SanitizeFileName(fileName), surfixName, extention))
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

}
