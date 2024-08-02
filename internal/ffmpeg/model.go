package ffmpeg

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type ffmpegModel struct {
	width  int
	height int

	stage int

	bubbles bubbleModel

	selectionOption SelectionOption

	questions Questions

	downloader downloader
}

type bubbleModel struct {
	spinner  spinner.Model
	progress progress.Model
}

type downloader struct {
	downloadPrecentage float64
	itag               int
	videoPath          string
	audioPath          string
	thumnailPath       string
	isProcessing       bool
	outputPath         string
}

type SelectionOption struct {
	wantToDownload bool
}

type Questions struct {
	question   string
	defaultVal string
	input      Input
}

type Input interface {
	Blink() tea.Msg
	View() string
	Update(msg tea.Msg) (Input, tea.Cmd)
	SetValue(msg string)
	Blur()
	Value() string
}

type answerField struct {
	textinput textinput.Model
}

func (a *answerField) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *answerField) View() string {
	return a.textinput.View()

}

func (a *answerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)

	return a, cmd
}

func (a *answerField) SetValue(msg string) {
	a.textinput.SetValue(msg)

}

func (a *answerField) Blur() {
	a.textinput.Blur()

}

func (a *answerField) Value() string {
	return a.textinput.Value()
}
