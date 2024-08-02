package ffmpeg

import (
	"strings"
)

func (m *ffmpegModel) View() string {

	var strBuilder strings.Builder

	switch m.stage {

	case 1:
		VSelectionStage(m, &strBuilder)
	case 2:
		VSelectionStage(m, &strBuilder)
		VConfigureStage(m, &strBuilder)

	}
	return strBuilder.String()
}
