package yt

import (
	"YTDownloaderCli/internal/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

						F, err := tea.LogToFile("format.json", "format")
						if err != nil {
							fmt.Println("Error opening log file:", err)
							return
						}

						m.data.video.Formats.Sort()
						// var y FormatList
						d, _ := json.Marshal(detail.Formats)

						F.WriteString("\n\n\n\n")
						F.WriteString(fmt.Sprintf("%v", string(d)))
						F.WriteString("\n\n\n\n")

						defer F.Close()

						utils.SetRequiredVideoFormat(detail)
						// F.WriteString("\n\n\n\n")
						// F.WriteString(fmt.Sprintf("%v", detail.Formats))
						// F.WriteString("\n\n\n\n")

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
	thumbnailPath := downloadThumnail(m)

	if strings.Contains(format.MimeType, "avc1") {

		videoPath := download(m, &format, m.data.video.Title, format.QualityLabel, "mp4")
		m.data.downloadPrecentage = float64(0)

		audioPath := download(m, utils.GetMaxAudioQuality(m.data.video.Formats), m.data.video.Title, format.AudioQuality, "m4a")
		mergeVideoAudio(m, videoPath, audioPath, thumbnailPath)

		defer os.RemoveAll(thumbnailPath)
		defer os.RemoveAll(videoPath)
		defer os.RemoveAll(audioPath)
	} else {
		download(m, utils.GetMaxAudioQuality(m.data.video.Formats), m.data.video.Title, format.AudioQuality, "m4a")
	}

}

func download(m *model, format *youtube.Format, fileName string, surfixName string, extention string) string {
	//* start downlod
	stream, s, err := m.client.GetDownloadStream(&m.data.video, format)

	utils.UtilError(err)

	defer stream.Close()

	file, err := os.Create(fmt.Sprintf("%s_%s.%s", utils.SanitizeFileName(fileName), surfixName, extention))
	utils.UtilError(err)

	defer file.Close()

	pw := &ProgressWriter{
		Writer:     file,
		TotalBytes: s,
		ProgressDisplay: func(progress float64) {
			m.data.downloadPrecentage = progress
		},
	}
	_, err = io.Copy(pw, stream)
	utils.UtilError(err)
	return file.Name()
}

func downloadThumnail(m *model) string {

	res, err := http.Get(m.data.video.Thumbnails[len(m.data.video.Thumbnails)-1].URL)

	thumnailName := fmt.Sprintf("%v", m.data.video.Title)
	utils.UtilError(err)
	defer res.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s.jpg", utils.SanitizeFileName(thumnailName)))

	io.Copy(file, res.Body)

	log := utils.Log()
	defer log.Close()
	log.WriteString(fmt.Sprintf("%v", m.data.video.Thumbnails[len(m.data.video.Thumbnails)-1]))

	return file.Name()
}

// mergeVideoAudio merges a video and audio file using FFmpeg
func mergeVideoAudio(m *model, videoPath, audioPath, thumbnailPath string) {

	ffmpegPath2, err := utils.ExtractFFmpeg()
	utils.UtilError(err)
	defer os.RemoveAll(filepath.Dir(ffmpegPath2)) // Clean up temporary directory

	// Open the log file
	F, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer F.Close()

	// Log paths
	F.WriteString(fmt.Sprintf("ffmpegPath: %v\n", ffmpegPath2))
	F.WriteString(fmt.Sprintf("videoPath: %v\n", videoPath))
	F.WriteString(fmt.Sprintf("audioPath: %v\n", audioPath))
	F.WriteString(fmt.Sprintf("thumbnailPath: %v\n", thumbnailPath))

	outputFileName := fmt.Sprintf("%v.mp4", m.data.video.Title)
	// Prepare FFmpeg command arguments
	args := []string{
		"-y",
		"-i", videoPath,
		"-i", thumbnailPath,
		"-i", audioPath,
		"-map", "0:v",
		"-map", "1",
		"-map", "2:a",
		"-c:v", "copy",
		"-c:a", "copy",
		"-c:v:1", "png",
		"-disposition:v:1", "attached_pic",
		outputFileName,
	}

	// Create and execute the command
	cmd := exec.Command(ffmpegPath2, args...)

	F.WriteString("\n\n")
	F.WriteString(fmt.Sprintf("Command: %v %v", ffmpegPath2, strings.Join(args, " ")))
	F.WriteString("\n\n")

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		F.WriteString(fmt.Sprintf("FFmpeg command failed: %v\nOutput: %s\n", err, output))
		return
	}

	F.WriteString("Video processing completed successfully.\n")
	F.WriteString(fmt.Sprintf("FFmpeg output: %s\n", output))
}
