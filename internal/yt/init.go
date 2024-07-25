package yt

import (
	ytclient "YTDownloaderCli/internal/service/ytClient"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

func newAnswerField() *answerField {
	a := answerField{}
	newInput := textinput.New()
	newInput.Placeholder = "https://youtu.be/l-BgjOr5FJY"
	newInput.Focus()
	a.textinput = newInput
	return &a
}

func newQuestion(question string) *Questions {

	q := Questions{question: question, defaultVal: "https://youtu.be/l-BgjOr5FJY"}

	model := newAnswerField()
	q.input = model
	return &q
}

func initialModel() *model {

	questions := []Questions{*newQuestion("provide any youtube video or playlist url")}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &model{
		client: *ytclient.NewYTClient(),

		questions: questions,
		index:     0,
		stage:     1,
		bubbles: bubbleModel{
			spinner:  s,
			progress: progress.New(progress.WithDefaultGradient()),
		},
		contentTypeSelection: SelectionModel{
			choices: []string{"Video", "Playlist"},
		},
	}
}

func (m model) Init() tea.Cmd {
	type tickMsg time.Time
	return tea.Batch(m.questions[m.index].input.Blink, m.bubbles.spinner.Tick,	tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	}))
}
