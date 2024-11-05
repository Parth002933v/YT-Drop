package utils

import (
	"embed"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	// tea "github.com/charmbracelet/bubbletea"

	"github.com/mitchellh/go-homedir"
)

func UtilError(erro error) {

	if erro != nil {
		fmt.Println("fatal:", erro)
		os.Exit(1)
	}
}

func HomeDir() string {
	home, err := homedir.Dir()
	UtilError(err)
	path := filepath.Join(home, ".yt-drop")
	return path
}

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

func BitseToMB(ContentLength int64) float64 {

	return float64(ContentLength) / (1024 * 1024)

}
func GetfprmatInFile(formats youtube.FormatList) {
	// Create a new file where the Go struct data will be written
	file, err := os.Create("formats_output.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write Go-formatted data to the file
	fmt.Fprintf(file, "package main\n\nvar formats = []youtube.Format{\n")
	for _, format := range formats {
		// Write each format as Go-struct formatted data
		escapedMimeType := strings.ReplaceAll(format.MimeType, `"`, `\"`)

		// Write each struct's data as Go-struct formatted text

		// Write each struct's data as Go-struct formatted text
		fmt.Fprintf(file, "  {    ItagNo: %d,  QualityLabel: \"%s\",   MimeType: \"%s\",    Quality: \"%s\",    Bitrate: %d,    FPS: %d,    Width: %d,    Height: %d,    LastModified: \"%s\",    ContentLength: %d,       ProjectionType: \"%s\",    AverageBitrate: %d,    AudioQuality: \"%s\",    ApproxDurationMs: \"%s\",    AudioSampleRate: \"%s\",    AudioChannels: %d, Cipher: \"%s\",  },\n",
			format.ItagNo, format.QualityLabel, escapedMimeType, format.Quality, format.Bitrate, format.FPS,
			format.Width, format.Height, format.LastModified, format.ContentLength,
			format.ProjectionType, format.AverageBitrate,
			format.AudioQuality, format.ApproxDurationMs, format.AudioSampleRate, format.AudioChannels,
			format.Cipher,
			//format.URL,
		)
	}
	fmt.Fprintf(file, "}\n")

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
func GetFormats() *youtube.FormatList {

	d := &youtube.FormatList{
		{ItagNo: 337, QualityLabel: "2160p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "hd2160", Bitrate: 29973590, FPS: 60, Width: 3840, Height: 2160, LastModified: "1711474575024876", ContentLength: 1135566526, ProjectionType: "RECTANGULAR", AverageBitrate: 28952008, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 315, QualityLabel: "2160p60", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd2160", Bitrate: 26703947, FPS: 60, Width: 3840, Height: 2160, LastModified: "1711474355309133", ContentLength: 1000488328, ProjectionType: "RECTANGULAR", AverageBitrate: 25508101, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 401, QualityLabel: "2160p", MimeType: "video/mp4; codecs=\"av01.0.12M.08\"", Quality: "hd2160", Bitrate: 11894322, FPS: 25, Width: 3840, Height: 2160, LastModified: "1693128858160744", ContentLength: 856148003, ProjectionType: "RECTANGULAR", AverageBitrate: 9322169, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 313, QualityLabel: "2160p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd2160", Bitrate: 17722981, FPS: 25, Width: 3840, Height: 2160, LastModified: "1693126537517537", ContentLength: 1585463695, ProjectionType: "RECTANGULAR", AverageBitrate: 17263324, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 336, QualityLabel: "1440p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "hd1440", Bitrate: 16551248, FPS: 60, Width: 2560, Height: 1440, LastModified: "1711467642314203", ContentLength: 637426396, ProjectionType: "RECTANGULAR", AverageBitrate: 16251601, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 308, QualityLabel: "1440p60", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd1440", Bitrate: 13299899, FPS: 60, Width: 2560, Height: 1440, LastModified: "1711466275273369", ContentLength: 408682929, ProjectionType: "RECTANGULAR", AverageBitrate: 10419637, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 400, QualityLabel: "1440p", MimeType: "video/mp4; codecs=\"av01.0.12M.08\"", Quality: "hd1440", Bitrate: 5345755, FPS: 25, Width: 2560, Height: 1440, LastModified: "1693126816005344", ContentLength: 354467607, ProjectionType: "RECTANGULAR", AverageBitrate: 3859621, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 271, QualityLabel: "1440p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd1440", Bitrate: 7270921, FPS: 25, Width: 2560, Height: 1440, LastModified: "1693126162626689", ContentLength: 461931271, ProjectionType: "RECTANGULAR", AverageBitrate: 5029739, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 335, QualityLabel: "1080p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "hd1080", Bitrate: 6963898, FPS: 60, Width: 1920, Height: 1080, LastModified: "1711466234708613", ContentLength: 263936467, ProjectionType: "RECTANGULAR", AverageBitrate: 6729232, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 303, QualityLabel: "1080p60", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd1080", Bitrate: 5796797, FPS: 60, Width: 1920, Height: 1080, LastModified: "1711469035327759", ContentLength: 134249528, ProjectionType: "RECTANGULAR", AverageBitrate: 3422779, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 399, QualityLabel: "1080p", MimeType: "video/mp4; codecs=\"av01.0.08M.08\"", Quality: "hd1080", Bitrate: 1681120, FPS: 25, Width: 1920, Height: 1080, LastModified: "1693124209186949", ContentLength: 111049412, ProjectionType: "RECTANGULAR", AverageBitrate: 1209161, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 248, QualityLabel: "1080p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd1080", Bitrate: 2439632, FPS: 25, Width: 1920, Height: 1080, LastModified: "1693131240266257", ContentLength: 125420271, ProjectionType: "RECTANGULAR", AverageBitrate: 1365638, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 137, QualityLabel: "1080p", MimeType: "video/mp4; codecs=\"avc1.640028\"", Quality: "hd1080", Bitrate: 2517352, FPS: 25, Width: 1920, Height: 1080, LastModified: "1693123075425138", ContentLength: 147825606, ProjectionType: "RECTANGULAR", AverageBitrate: 1609599, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 334, QualityLabel: "720p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "hd720", Bitrate: 4540895, FPS: 60, Width: 1280, Height: 720, LastModified: "1711462402821837", ContentLength: 169907732, ProjectionType: "RECTANGULAR", AverageBitrate: 4331908, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 302, QualityLabel: "720p60", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd720", Bitrate: 3470975, FPS: 60, Width: 1280, Height: 720, LastModified: "1711462418099871", ContentLength: 74193960, ProjectionType: "RECTANGULAR", AverageBitrate: 1891623, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 398, QualityLabel: "720p", MimeType: "video/mp4; codecs=\"av01.0.05M.08\"", Quality: "hd720", Bitrate: 933082, FPS: 25, Width: 1280, Height: 720, LastModified: "1693123260288823", ContentLength: 56787290, ProjectionType: "RECTANGULAR", AverageBitrate: 618328, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 247, QualityLabel: "720p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "hd720", Bitrate: 1249258, FPS: 25, Width: 1280, Height: 720, LastModified: "1693124539671231", ContentLength: 65312660, ProjectionType: "RECTANGULAR", AverageBitrate: 711157, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 136, QualityLabel: "720p", MimeType: "video/mp4; codecs=\"avc1.4d401f\"", Quality: "hd720", Bitrate: 1083228, FPS: 25, Width: 1280, Height: 720, LastModified: "1693123628951416", ContentLength: 44196412, ProjectionType: "RECTANGULAR", AverageBitrate: 481232, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 333, QualityLabel: "480p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "large", Bitrate: 1989702, FPS: 60, Width: 854, Height: 480, LastModified: "1711462404421641", ContentLength: 72263839, ProjectionType: "RECTANGULAR", AverageBitrate: 1842413, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 397, QualityLabel: "480p", MimeType: "video/mp4; codecs=\"av01.0.04M.08\"", Quality: "large", Bitrate: 475562, FPS: 25, Width: 854, Height: 480, LastModified: "1693125734421494", ContentLength: 28110538, ProjectionType: "RECTANGULAR", AverageBitrate: 306081, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 244, QualityLabel: "480p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "large", Bitrate: 709661, FPS: 25, Width: 854, Height: 480, LastModified: "1693124467306256", ContentLength: 35495810, ProjectionType: "RECTANGULAR", AverageBitrate: 386496, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 135, QualityLabel: "480p", MimeType: "video/mp4; codecs=\"avc1.4d401e\"", Quality: "large", Bitrate: 511792, FPS: 25, Width: 854, Height: 480, LastModified: "1693123309496928", ContentLength: 25400361, ProjectionType: "RECTANGULAR", AverageBitrate: 276571, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 332, QualityLabel: "360p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "medium", Bitrate: 1062858, FPS: 60, Width: 640, Height: 360, LastModified: "1711466218221965", ContentLength: 37439972, ProjectionType: "RECTANGULAR", AverageBitrate: 954556, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 396, QualityLabel: "360p", MimeType: "video/mp4; codecs=\"av01.0.01M.08\"", Quality: "medium", Bitrate: 272242, FPS: 25, Width: 640, Height: 360, LastModified: "1693125528175318", ContentLength: 15869634, ProjectionType: "RECTANGULAR", AverageBitrate: 172796, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 243, QualityLabel: "360p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "medium", Bitrate: 434461, FPS: 25, Width: 640, Height: 360, LastModified: "1693124663346082", ContentLength: 23022312, ProjectionType: "RECTANGULAR", AverageBitrate: 250678, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		//{ItagNo: 18, QualityLabel: "360p", MimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"", Quality: "medium", Bitrate: 288711, FPS: 25, Width: 640, Height: 360, LastModified: "1699199506995285", ContentLength: 0, ProjectionType: "RECTANGULAR", AverageBitrate: 0, AudioQuality: "AUDIO_QUALITY_LOW", ApproxDurationMs: "734772", AudioSampleRate: "44100", AudioChannels: 2, Cipher: ""},
		{ItagNo: 134, QualityLabel: "360p", MimeType: "video/mp4; codecs=\"avc1.4d401e\"", Quality: "medium", Bitrate: 282138, FPS: 25, Width: 640, Height: 360, LastModified: "1693123066134796", ContentLength: 14655061, ProjectionType: "RECTANGULAR", AverageBitrate: 159571, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 331, QualityLabel: "240p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "small", Bitrate: 500134, FPS: 60, Width: 426, Height: 240, LastModified: "1711462437710721", ContentLength: 17370232, ProjectionType: "RECTANGULAR", AverageBitrate: 442865, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 395, QualityLabel: "240p", MimeType: "video/mp4; codecs=\"av01.0.00M.08\"", Quality: "small", Bitrate: 132776, FPS: 25, Width: 426, Height: 240, LastModified: "1693125299444869", ContentLength: 8601916, ProjectionType: "RECTANGULAR", AverageBitrate: 93661, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 242, QualityLabel: "240p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "small", Bitrate: 207936, FPS: 25, Width: 426, Height: 240, LastModified: "1693124773573198", ContentLength: 10172613, ProjectionType: "RECTANGULAR", AverageBitrate: 110764, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 133, QualityLabel: "240p", MimeType: "video/mp4; codecs=\"avc1.4d4015\"", Quality: "small", Bitrate: 157384, FPS: 25, Width: 426, Height: 240, LastModified: "1693123196046218", ContentLength: 7860422, ProjectionType: "RECTANGULAR", AverageBitrate: 85588, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 330, QualityLabel: "144p60 HDR", MimeType: "video/webm; codecs=\"vp9.2\"", Quality: "tiny", Bitrate: 244489, FPS: 60, Width: 256, Height: 144, LastModified: "1711462411714910", ContentLength: 8297297, ProjectionType: "RECTANGULAR", AverageBitrate: 211544, AudioQuality: "", ApproxDurationMs: "313779", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 394, QualityLabel: "144p", MimeType: "video/mp4; codecs=\"av01.0.00M.08\"", Quality: "tiny", Bitrate: 74504, FPS: 25, Width: 256, Height: 144, LastModified: "1693125177593465", ContentLength: 5034082, ProjectionType: "RECTANGULAR", AverageBitrate: 54813, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 278, QualityLabel: "144p", MimeType: "video/webm; codecs=\"vp9\"", Quality: "tiny", Bitrate: 107727, FPS: 25, Width: 256, Height: 144, LastModified: "1693124937583866", ContentLength: 6690173, ProjectionType: "RECTANGULAR", AverageBitrate: 72845, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 160, QualityLabel: "144p", MimeType: "video/mp4; codecs=\"avc1.4d400c\"", Quality: "tiny", Bitrate: 77202, FPS: 25, Width: 256, Height: 144, LastModified: "1693123157953991", ContentLength: 3936421, ProjectionType: "RECTANGULAR", AverageBitrate: 42861, AudioQuality: "", ApproxDurationMs: "734720", AudioSampleRate: "", AudioChannels: 0, Cipher: ""},
		{ItagNo: 140, QualityLabel: "", MimeType: "audio/mp4; codecs=\"mp4a.40.2\"", Quality: "tiny", Bitrate: 130982, FPS: 0, Width: 0, Height: 0, LastModified: "1693122990257594", ContentLength: 11892179, ProjectionType: "RECTANGULAR", AverageBitrate: 129478, AudioQuality: "AUDIO_QUALITY_MEDIUM", ApproxDurationMs: "734772", AudioSampleRate: "44100", AudioChannels: 2, Cipher: ""},
		{ItagNo: 139, QualityLabel: "", MimeType: "audio/mp4; codecs=\"mp4a.40.5\"", Quality: "tiny", Bitrate: 50537, FPS: 0, Width: 0, Height: 0, LastModified: "1693122908403968", ContentLength: 4481743, ProjectionType: "RECTANGULAR", AverageBitrate: 48789, AudioQuality: "AUDIO_QUALITY_LOW", ApproxDurationMs: "734865", AudioSampleRate: "22050", AudioChannels: 2, Cipher: ""},
		{ItagNo: 251, QualityLabel: "", MimeType: "audio/webm; codecs=\"opus\"", Quality: "tiny", Bitrate: 139425, FPS: 0, Width: 0, Height: 0, LastModified: "1693122907906711", ContentLength: 11053452, ProjectionType: "RECTANGULAR", AverageBitrate: 120352, AudioQuality: "AUDIO_QUALITY_MEDIUM", ApproxDurationMs: "734741", AudioSampleRate: "48000", AudioChannels: 2, Cipher: ""},
		{ItagNo: 250, QualityLabel: "", MimeType: "audio/webm; codecs=\"opus\"", Quality: "tiny", Bitrate: 77378, FPS: 0, Width: 0, Height: 0, LastModified: "1693122827171660", ContentLength: 6039065, ProjectionType: "RECTANGULAR", AverageBitrate: 65754, AudioQuality: "AUDIO_QUALITY_LOW", ApproxDurationMs: "734741", AudioSampleRate: "48000", AudioChannels: 2, Cipher: ""},
		{ItagNo: 249, QualityLabel: "", MimeType: "audio/webm; codecs=\"opus\"", Quality: "tiny", Bitrate: 54990, FPS: 0, Width: 0, Height: 0, LastModified: "1693122912505033", ContentLength: 4699689, ProjectionType: "RECTANGULAR", AverageBitrate: 51171, AudioQuality: "AUDIO_QUALITY_LOW", ApproxDurationMs: "734741", AudioSampleRate: "48000", AudioChannels: 2, Cipher: ""},
	}

	return d

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

func GetVTypeAndCodecFromMimType(str string) (Vtype MediaType, Codec string, error error) {
	// Regular expression to capture the codec identifier (e.g., avc1, av01, vp9)
	re := regexp.MustCompile(`(video|audio)/[a-z0-9]+; codecs="([a-zA-Z0-9]+)`)

	match := re.FindStringSubmatch(str)
	if len(match) > 2 {
		return GetMediaType(match[1]), match[2], nil
	} else {
		return GetMediaType("Unknown Codec"), "Unknown Codec", fmt.Errorf("unknown mimetype")
	}
}

func TruncateWithEllipsisString(s string, maxLength int) string {
	if len(s) > maxLength {
		return fmt.Sprintf("%s...", s[:maxLength-3])
	}
	return s
}
func LogToFileWith(path string, prefix string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600) //nolint:gomnd
	if err != nil {
		return nil, fmt.Errorf("error opening file for logging: %w", err)
	}

	// Add a space after the prefix if a prefix is being specified, and it
	// doesn't already have a trailing space.
	if len(prefix) > 0 {
		finalChar := prefix[len(prefix)-1]
		if !unicode.IsSpace(rune(finalChar)) {
			prefix += " "
		}
	}

	return f, nil
}
