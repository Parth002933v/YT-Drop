package yt

import (
	"strings"
)

func (m *model) View() string {
	var strBuilder strings.Builder

	text := "YT Downloader"
	width := len(text) + 4

	// Draw the border and header
	border := strings.Repeat("-", width)
	strBuilder.WriteString("\n")
	strBuilder.WriteString(border)
	strBuilder.WriteString("\n| ")
	strBuilder.WriteString(text)
	strBuilder.WriteString(" |\n")
	strBuilder.WriteString(border)

	strBuilder.WriteString("\n\n")

	switch m.stage {

	case 1:
		VSelectionStage(m, &strBuilder)

	case 2:
		VSelectionStage(m, &strBuilder)
		VUrlInputStage(m, &strBuilder)

	case 3:
		VSelectionStage(m, &strBuilder)
		VUrlInputStage(m, &strBuilder)
		VResposeDataStage(m, &strBuilder)

	case 4:
		VSelectionStage(m, &strBuilder)
		VUrlInputStage(m, &strBuilder)
		VResposeDataStage(m, &strBuilder)
		VDownloadProgressStage(m, &strBuilder)
	}

	strBuilder.WriteString("\n")

	finalWidth := len(strBuilder.String())
	finalBorder := strings.Repeat("-", finalWidth)

	return finalBorder + "\n" + strBuilder.String() + finalBorder

}
