package ffmpeg

import (
	utils "YTDownloaderCli/internal/common"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func VSelectionStage(m *ffmpegModel, strBuilder *strings.Builder) {
	redFont := lipgloss.NewStyle().Italic(true).Bold(true).Foreground(lipgloss.Color("#FF0000"))
	italic := lipgloss.NewStyle().Italic(true).Bold(true)

	strBuilder.WriteString("\n")

	strBuilder.WriteString(redFont.Render("The application is not able to locate the ffmpeg"))

	strBuilder.WriteString("\n")
	strBuilder.WriteString("-Press ")
	strBuilder.WriteString(italic.Render("locate"))
	strBuilder.WriteString(" if you already have FFmpeg")
	strBuilder.WriteString("\n")

	strBuilder.WriteString("-Press ")
	strBuilder.WriteString(italic.Render("download"))
	strBuilder.WriteString(" to download and configure FFmpeg automatically")

	strBuilder.WriteString("\n\n")

	hotPink := lipgloss.Color("#FF06B7")
	Selectedstyle := lipgloss.NewStyle().Foreground(hotPink)

	if m.selectionOption.wantToDownload {
		strBuilder.WriteString("Locate")
		strBuilder.WriteString("\t\t")
		strBuilder.WriteString(Selectedstyle.Render("download"))

	} else {
		strBuilder.WriteString(Selectedstyle.Render("Locate"))
		strBuilder.WriteString("\t\t")
		strBuilder.WriteString("download")
	}

	strBuilder.WriteString("\n\n")
}

func VConfigureStage(m *ffmpegModel, strBuilder *strings.Builder) {
	if m.selectionOption.wantToDownload {
		downloadAndConfigureFFMpeg(m, strBuilder)
	} else {
		specifyFFMpegPath(m, strBuilder)
	}
}

func specifyFFMpegPath(m *ffmpegModel, strBuilder *strings.Builder) {

	var strInBoxBuilder strings.Builder
	box := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).PaddingLeft(2).PaddingRight(2)

	strInBoxBuilder.WriteString("Please Provide The exact path where you installed your FFMpeg\n")
	strInBoxBuilder.WriteString(utils.WrapText(m.questions.input.View(), m.width))

	strBuilder.WriteString(box.Render(strInBoxBuilder.String()))
}

func downloadAndConfigureFFMpeg(m *ffmpegModel, strBuilder *strings.Builder) {
	if m.downloader.isProcessing {
		strBuilder.WriteString("Preparing to Download FFMpeg...")

		strBuilder.WriteString("\n")

		strBuilder.WriteString(m.bubbles.progress.ViewAs(m.downloader.downloadPrecentage))

		strBuilder.WriteString("\n")

		strBuilder.WriteString("Preparing to Download FFMpeg...")

	}
}
