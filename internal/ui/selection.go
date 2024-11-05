package ui

import (
	"YTDownloaderCli/internal/sharedState"
	"YTDownloaderCli/internal/utils"
	"YTDownloaderCli/pkg/_youtube"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/pterm/pterm"
)

func TypeSelection() (sharedState.DownloadType, error) {
	var error error
	options := []string{sharedState.TypeVideo.String(), sharedState.TypePlaylist.String()}

	ptermSelection := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithFilter(false).
		WithDefaultText("Choose what type of content you want to download")

	ptermSelection.OptionStyle.Code()
	ptermSelection.TextStyle = pterm.NewStyle(pterm.FgGray)
	selectedOption, selectedOptionErr := ptermSelection.Show()

	if selectedOptionErr != nil {
		error = selectedOptionErr
	}

	downloadType, err := sharedState.DownloadTypeFromString(selectedOption)
	if err != nil {
		error = err
	}
	return downloadType, error
}

func FormatSelection(formats []youtube.Format) youtube.Format {
	maxAudio := utils.GetMaxAudioQuality(formats)

	var formatedOptions []string
	for _, format := range formats {
		mediaType, codec, _ := utils.GetVTypeAndCodecFromMimType(format.MimeType)
		if mediaType == utils.Video {
			option := fmt.Sprintf("Video: | %s | %.2f MB | %s", format.QualityLabel, utils.BitseToMB(format.ContentLength)+utils.BitseToMB(maxAudio.ContentLength), codec)
			formatedOptions = append(formatedOptions, option)
		} else if mediaType == utils.Audio {
			option := fmt.Sprintf("Audio: | %s | %.2f MB | %s", format.AudioSampleRate, utils.BitseToMB(format.ContentLength)+utils.BitseToMB(maxAudio.ContentLength), codec)
			formatedOptions = append(formatedOptions, option)
		} else {
			formatedOptions = append(formatedOptions, "Unknown format")
		}
	}

	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithOptions(formatedOptions).
		Show("Select download format:")

	var selectedFormat youtube.Format
	for i, option := range formatedOptions {
		if option == selectedOption {
			selectedFormat = formats[i]
			break
		}
	}

	return selectedFormat

}

func MultiVideoSelection(videos []*_youtube.PlaylistEntry) []*_youtube.PlaylistEntry {
	var fomatedVideos []string

	for i, video := range videos {
		option := fmt.Sprintf("%d. %s", i+1, video.Title)
		fomatedVideos = append(fomatedVideos, option)
	}
	// Create a new interactive multiselect printer with the options
	// Disable the filter and define the checkmark symbols
	selectedOption, _ := pterm.DefaultInteractiveMultiselect.
		WithOptions(fomatedVideos).
		WithFilter(true).
		WithDefaultText("Select Video You want to download").
		WithCheckmark(&pterm.Checkmark{Checked: pterm.Green("✓"), Unchecked: pterm.Red(" ")}).Show()

	var selectedVideos []*_youtube.PlaylistEntry

	for i, option := range fomatedVideos {
		for _, selectedOption := range selectedOption {
			if option == selectedOption {
				selectedVideos = append(selectedVideos, videos[i])
				break
			}
		}

	}

	return selectedVideos
}
