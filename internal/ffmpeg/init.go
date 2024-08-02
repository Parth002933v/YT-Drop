package ffmpeg

import (
	utils "YTDownloaderCli/internal/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func newAnswerField(placeholder string) *answerField {
	a := answerField{}
	newInput := textinput.New()
	newInput.Placeholder = placeholder
	newInput.Focus()
	a.textinput = newInput
	return &a
}

func newQuestion(question string, placeholder string) *Questions {

	q := Questions{question: question}

	model := newAnswerField(placeholder)
	q.input = model
	return &q
}

func initialModel() *ffmpegModel {

	question := *newQuestion("Provide full path of your ffmpeg path", utils.HomeDir())

	return &ffmpegModel{stage: 1,
		questions: question,
	}
}

func (m *ffmpegModel) Init() tea.Cmd {
	return tea.Cmd(m.questions.input.Blink)
}
