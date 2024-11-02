package yt

import (
	"YTDownloaderCli/pkg/_youtube"
	"io"
	"sync"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kkdai/youtube/v2"
)

type model struct {
	client _youtube.YTClientModel

	width  int
	height int

	isPlaylist bool

	stage int8

	index                int
	questions            []Questions
	contentTypeSelection SelectionModel
	QualitySelection     SelectionModel

	bubbles bubbleModel

	downloader []downloader

	isLoading bool
	data      VideoResData
	done      bool
}

type downloader struct {
	itag         int
	downloaded   float64
	videoPath    string
	audioPath    string
	thumnailPath string
	isProcessing bool
	outputPath   string
}

type bubbleModel struct {
	spinner  spinner.Model
	progress progress.Model
}

type SelectionModel struct {
	choices []string
	cursor  int
}

type Questions struct {
	question   string
	defaultVal string
	input      Input
}

type VideoResData struct {
	video              youtube.Video
	playlist           youtube.Playlist
	downloadPrecentage float64
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

type ProgressWriter struct {
	Writer          io.Writer
	TotalBytes      int64
	BytesWritten    int64
	ProgressDisplay func(float64)
	mu              sync.Mutex
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	n, err = pw.Writer.Write(p)
	pw.BytesWritten += int64(n)

	if pw.ProgressDisplay != nil {
		progress := float64(pw.BytesWritten) / float64(pw.TotalBytes)
		pw.ProgressDisplay(progress)
	}

	return n, err
}
