package yt

import (
	"YTDownloaderCli/internal/common"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

						// F, err := tea.LogToFile("format.json", "format")
						// if err != nil {
						// 	fmt.Println("Error opening log file:", err)
						// 	return
						// }

						m.data.video.Formats.Sort()
						// var y FormatList
						// d, _ := json.Marshal(detail.Formats)

						// F.WriteString("\n\n\n\n")
						// F.WriteString(fmt.Sprintf("%v", string(d)))
						// F.WriteString("\n\n\n\n")

						// defer F.Close()

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
	chapterspath := extractChaptersAsFile(m)

	if strings.Contains(format.MimeType, "avc1") {

		videoPath := download(m, &format, m.data.video.Title, format.QualityLabel, "mp4")
		m.data.downloadPrecentage = float64(0)

		audioPath := download(m, utils.GetMaxAudioQuality(m.data.video.Formats), m.data.video.Title, format.AudioQuality, "m4a")

		if len(chapterspath) > 0 {
			mergeVideoAudioThumbnailChapters(m, videoPath, audioPath, thumbnailPath, &chapterspath)
		} else {
			mergeVideoAudioThumbnailChapters(m, videoPath, audioPath, thumbnailPath, nil)
		}

		defer os.RemoveAll(thumbnailPath)
		defer os.RemoveAll(videoPath)
		defer os.RemoveAll(audioPath)
		defer os.RemoveAll(chapterspath)
	} else {
		download(m, utils.GetMaxAudioQuality(m.data.video.Formats), m.data.video.Title, format.AudioQuality, "m4a")
	}

}

func download(m *model, format *youtube.Format, fileName string, surfixName string, extention string) string {
	//* start downlod
	stream, s, err := m.client.GetDownloadStreamWithContext(&m.data.video, format)

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
	utils.UtilError(err)
	defer file.Close()

	io.Copy(file, res.Body)

	// log := utils.Log()
	// defer log.Close()
	// log.WriteString(fmt.Sprintf("%v", m.data.video.Thumbnails[len(m.data.video.Thumbnails)-1]))

	return file.Name()
}

func extractChapters(description string) []string {
	re := regexp.MustCompile(`\d{2}:\d{2} .+`)
	return re.FindAllString(description, -1)
}

func saveChaptersToFile(chapters []string, videoTitle string) string {
	fileName := fmt.Sprintf("%s_chapters.txt", utils.SanitizeFileName(videoTitle))
	filePath := filepath.Join(".", fileName)

	file, err := os.Create(filePath)
	utils.UtilError(err)

	defer file.Close()

	content := strings.Join(chapters, "\n")
	_, err = file.WriteString(content)
	utils.UtilError(err)
	return filePath
}

func convertChaterToFFmpegFormat() {}
func extractChaptersAsFile(m *model) string {

	chapters := extractChapters(m.data.video.Description)
	if len(chapters) > 0 {
		chaptersFilePath := saveChaptersToFile(chapters, m.data.video.Title)
		return chaptersFilePath
	} else {
		return ""
	}
}

// mergeVideoAudioThumbnailChapters merges a video and audio file using FFmpeg
func mergeVideoAudioThumbnailChapters(m *model, videoPath, audioPath, thumbnailPath string, chaptersPath *string) {

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
	F.WriteString(fmt.Sprintf("chapterPath: %v\n", thumbnailPath))

	outputFileName := fmt.Sprintf("%v.mp4", utils.SanitizeFileName(m.data.video.Title))
	// // Prepare FFmpeg command arguments
	// args := []string{
	// 	"-y",
	// 	"-i", videoPath,
	// 	"-i", thumbnailPath,
	// 	"-i", audioPath,
	// 	"-map", "0:v",
	// 	"-map", "1",
	// 	"-map", "2:a",
	// 	"-c:v", "copy",
	// 	"-c:a", "copy",
	// 	"-c:v:1", "png",
	// 	"-disposition:v:1", "attached_pic",
	// 	outputFileName,
	// }

	args := []string{
		"-y",
		"-i", videoPath,
		"-i", thumbnailPath,
		"-i", audioPath,
	}

	// Conditionally add chaptersPath if it's not nil
	// if chaptersPath != nil && *chaptersPath != "" {
	// 	args = append(args, "-i", *chaptersPath)
	// }

	// Add the remaining arguments
	args = append(args,
		"-map", "0:v",
		"-map", "1",
		"-map", "2:a",
		"-c:v", "copy",
		"-c:a", "copy",
		"-c:v:1", "png",
		"-disposition:v:1", "attached_pic",
	)

	// If chaptersPath was added, map metadata
	// if chaptersPath != nil && *chaptersPath != "" {
	// 	args = append(args, "-map_metadata", fmt.Sprintf("%d", len(args)-1))
	// }

	args = append(args, outputFileName)

	// Create and execute the command
	cmd := exec.Command(ffmpegPath2, args...)

	F.WriteString("\n\n")
	F.WriteString(fmt.Sprintf("Command: %v %v", ffmpegPath2, strings.Join(args, " ")))
	F.WriteString("\n\n")

	// Capture output
	_, err = cmd.CombinedOutput()
	if err != nil {
		F.WriteString(fmt.Sprintf("FFmpeg command failed: %v", err))
		return
	}

	// F.WriteString("Video processing completed successfully.\n")
	// F.WriteString(fmt.Sprintf("FFmpeg output: %s\n", output))
}
