package yt

import (
	"YTDownloaderCli/internal/common"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kkdai/youtube/v2"
)

func VSelectionStage(m *model, strBuilder *strings.Builder) {

	strBuilder.WriteString("Choose what type of content you want to download")
	strBuilder.WriteString("\n")

	for i, choice := range m.contentTypeSelection.choices {
		cursor := " "
		isSelected := false
		if m.contentTypeSelection.cursor == i {
			cursor = ">"
			isSelected = true
		}

		selectedTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Bold(true)

		if m.stage > 1 {
			if isSelected {
				strBuilder.WriteString(fmt.Sprintf("%s %s\n", cursor, selectedTextStyle.Render(choice)))
				strBuilder.WriteString("\n")
			}
		} else {

			if isSelected {
				strBuilder.WriteString(fmt.Sprintf("%s %s\n", cursor, selectedTextStyle.Render(choice)))
			} else {

				strBuilder.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
			}
		}

	}
}

func VUrlInputStage(m *model, strBuilder *strings.Builder) 	{
	strBuilder.WriteString(utils.WrapText(m.questions[m.index].input.View(), m.width))

	strBuilder.WriteString("\n")
	if m.isLoading && m.data.video.ID == "" && m.data.playlist.ID == "" {
		strBuilder.WriteString(m.bubbles.spinner.View())
		strBuilder.WriteString("Loading...\n\n")
	}

}

func VResposeDataStage(m *model, strBuilder *strings.Builder) {
	if m.data.video.ID != "" || m.data.playlist.ID == "" {
		buildVideoMetadata(m, strBuilder)
		buildFormatSelectionBox(m, strBuilder)
	}
}

func VDownloadProgressStage(m *model, strBuilder *strings.Builder) {

	strBuilder.WriteString("\n")

	strBuilder.WriteString(m.bubbles.progress.ViewAs(m.data.downloadPrecentage))

}

// *=========================================================
func buildVideoMetadata(m *model, strBuilder *strings.Builder) {
	var metadataBuilder strings.Builder
	metadataBuilder.WriteString(fmt.Sprintf("Author   : %s\n", m.data.video.Author))
	metadataBuilder.WriteString(fmt.Sprintf("Title    : %s\n", m.data.video.Title))
	metadataBuilder.WriteString(fmt.Sprintf("Duration : %s\n", m.data.video.Duration))
	metadataBuilder.WriteString(fmt.Sprintf("Views    : %v", m.data.video.Views))

	metadataBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		PaddingLeft(1).PaddingRight(1).
		Render(metadataBuilder.String())

	strBuilder.WriteString(metadataBox)
	strBuilder.WriteString("\n")
}

func buildFormatSelectionBox(m *model, strBuilder *strings.Builder) {
	var formatBuilder strings.Builder
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("201")).
		PaddingLeft(1).PaddingRight(1)

	m.data.video.Formats.Sort()
	maxAudioFormat := utils.GetMaxAudioQuality(m.data.video.Formats)

	if m.stage > 3 {
		cursor, isSelected := getCursorAndSelectionStatus(m, m.QualitySelection.cursor)
		formatDetails := getFormatDetails(m.data.video.Formats[m.QualitySelection.cursor], *maxAudioFormat, cursor, isSelected)
		formatBuilder.WriteString(formatDetails)
		formatBuilder.WriteString("\n")
	} else {
		for i, format := range m.data.video.Formats {
			cursor, isSelected := getCursorAndSelectionStatus(m, i)
			formatDetails := getFormatDetails(format, *maxAudioFormat, cursor, isSelected)

			formatBuilder.WriteString(formatDetails)
			formatBuilder.WriteString("\n")
		}
	}

	strBuilder.WriteString(borderStyle.Render(formatBuilder.String()))
}

func getCursorAndSelectionStatus(m *model, index int) (string, bool) {
	if m.QualitySelection.cursor == index {
		return ">", true
	}
	return " ", false
}

func getFormatDetails(format, maxAudioFormat youtube.Format, cursor string, isSelected bool) string {

	var formatDetails string
	style := lipgloss.NewStyle()
	if isSelected {
		style = style.Foreground(lipgloss.Color("#FAFAFA")).Bold(true)
	}

	if strings.Contains(format.MimeType, "avc1") {
		videoLengthMB := float64(format.ContentLength) / (1024 * 1024)
		audioLengthMB := float64(maxAudioFormat.ContentLength) / (1024 * 1024)
		totalLengthMB := videoLengthMB + audioLengthMB

		formatDetails = fmt.Sprintf("%s Video: %v | %.2f MB", cursor, format.QualityLabel, totalLengthMB)
	} else if strings.Contains(format.MimeType, "mp4a") {
		audioLengthMB := float64(format.ContentLength) / (1024 * 1024)

		formatDetails = fmt.Sprintf("%s audio: %v | %.2f MB | %v", cursor, format.AudioSampleRate, audioLengthMB, format.AudioQuality)
	}

	return style.Render(formatDetails)
}

//*=========================================================
