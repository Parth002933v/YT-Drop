package utils

import (
	"YTDownloaderCli/internal/ui/theme"
	"embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	// tea "github.com/charmbracelet/bubbletea"

	"github.com/kkdai/youtube/v2"
	"github.com/mitchellh/go-homedir"
)

func UtilError(erro error) {

	if erro != nil {
		fmt.Println()
		theme.PrintErrorText(fmt.Sprintf("fatal: %v", erro))
		fmt.Println()
		os.Exit(1)
	}
}

func DetermineYouTubeUrlType(u string) string {
	parsedURL, err := url.Parse(u)
	if err != nil {
		UtilError(err)
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")
	playlistID := queryParams.Get("list")

	if videoID != "" && playlistID == "" {
		return "video"
	} else if playlistID != "" && videoID == "" {
		return "playlist"
	} else if videoID != "" && playlistID != "" {
		return "video within playlist"
	} else {
		return "unknown"
	}
}

func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var wrapped strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		for len(line) > width {
			wrapped.WriteString(line[:width])
			wrapped.WriteString("\n")
			line = line[width:]
		}
		wrapped.WriteString(line)
		wrapped.WriteString("\n")
	}
	return wrapped.String()
}

func ConvAudioSampleRateToKbps(AudioSampleRate string) float64 {

	sampleRateInt, _ := strconv.Atoi(AudioSampleRate)
	bitDepth := 16
	channels := 2
	// PError(err)

	bitrate := float64(sampleRateInt*bitDepth*channels) / 1000

	return bitrate

}

func ConvertToMB(size int64) float64 {
	// contentLength := 78867345
	contentLength := 78867345

	// Convert bytes to megabytes
	contentLengthMB := float64(contentLength) / (1024 * 1024)
	return contentLengthMB

}

// get only selected fomated which is good
func SetRequiredVideoFormat(v *youtube.Video) {

	var format youtube.FormatList
	for i, val := range v.Formats {
		fmt.Print(i)
		if (strings.Contains(val.MimeType, "avc1") || strings.Contains(val.MimeType, "mp4a")) && val.ItagNo != 18 {
			format = append(format, val)
		}
	}
	v.Formats = format
}

// set the formate list values in m.QualitySelection.choices
//
// in simple word add the formate list in multiselection option
func SetQualitySelectionChoiceValue(list *[]string, data youtube.FormatList) {
	var values []string
	for _, val := range data {
		values = append(values, strconv.Itoa(val.ItagNo))
	}

	*list = values
}

func SanitizeFileName(fileName string) string {
	var builder strings.Builder

	for _, r := range fileName {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			builder.WriteRune(r)
		}
	}

	// Convert to a valid string by trimming extra spaces and ensuring consistent spaces
	return strings.Join(strings.Fields(builder.String()), " ")
}

func ExecutableDir() string {

	path, _ := os.Executable()

	return filepath.Dir(path)
}

func HomeDir() string {
	home, err := homedir.Dir()
	UtilError(err)
	path := filepath.Join(home, ".ytd")
	return path
}

// func Log() *os.File {

// 	F, err := tea.LogToFile("debug.log", "debug")
// 	if err != nil {
// 		fmt.Println("Error opening log file:", err)
// 		return nil
// 	}

// 	return F
// 	// defer F.Close()
// }
////////////////////=================================================

//go:embed ffmpeg.exe
var FFmpegFS embed.FS

func ExtractFFmpeg() (string, error) {
	tempDir, err := os.MkdirTemp("", "ffmpeg") //C:\Users\pp542\AppData\Local\Temp\ffmpeg74917441
	if err != nil {
		return "", err
	}

	ffmpegPath := filepath.Join(tempDir, "ffmpeg.exe") // C:\Users\pp542\AppData\Local\Temp\ffmpeg74917441\ffmpeg.exe
	ffmpegData, err := FFmpegFS.ReadFile("ffmpeg.exe")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(ffmpegPath, ffmpegData, 0755)
	if err != nil {
		return "", err
	}

	// // Open the log file
	// F, err := tea.LogToFile("ffmpeg.log", "ffmpeg")
	// F.WriteString(fmt.Sprintf("tempDir: %v\n", tempDir))
	// F.WriteString(fmt.Sprintf("ffmpegPath: %v\n", ffmpegPath))
	// defer F.Close()

	return ffmpegPath, nil
}

func FilterFormatsByMineType(formatList *youtube.FormatList, mimetypes ...string) {

	var newFormatList []youtube.Format

	for _, format := range *formatList {
		for _, mineType := range mimetypes {
			if strings.Contains(format.MimeType, mineType) {
				newFormatList = append(newFormatList, format)
				break
			}
		}
	}

	*formatList = newFormatList
}
func GetMaxAudioQuality(f youtube.FormatList) *youtube.Format {
	var format youtube.Format

	for _, val := range f {
		if strings.Contains(val.AudioQuality, "AUDIO_QUALITY_MEDIUM") && strings.Contains(val.MimeType, "opus") {
			format = val
		}
	}
	return &format
}
