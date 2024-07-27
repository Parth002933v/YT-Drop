package utils

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/kkdai/youtube/v2"
)

func PError(erro error) {

	if erro != nil {
		fmt.Println("fatal:", erro)
		os.Exit(1)
	}
}

func DetermineYouTubeUrlType(u string) string {
	parsedURL, err := url.Parse(u)
	if err != nil {
		PError(err)
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

func GetMaxAudioQuality(f youtube.FormatList) *youtube.Format {
	var format youtube.Format

	for _, val := range f {
		if strings.Contains(val.MimeType, "mp4a") && val.ItagNo != 18 && val.AudioQuality == "AUDIO_QUALITY_MEDIUM" {
			format = val
		}
	}
	return &format
}

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
