package ui

import (
	"YTDownloaderCli/internal/sharedState"
	"YTDownloaderCli/internal/utils"
	"errors"

	"github.com/pterm/pterm"
)

func UrlInput(downloadType sharedState.DownloadType) (string, error) {

	newInput := pterm.DefaultInteractiveTextInput.WithDefaultText("URL")

	if downloadType == sharedState.TypeVideo {
		newInput = newInput.WithDefaultValue("https://youtu.be/LXb3EKWsInQ")
	} else if downloadType == sharedState.TypePlaylist {
		newInput = newInput.WithDefaultValue("https://youtube.com/playlist?list=PLbtI3_MArDOk7J-8hR6CeB5U6bvgRKNNr")
	} else {
		utils.UtilError(errors.New("in correct download type"))
	}

	makeInput, err := newInput.Show()

	// result, err := pterm.DefaultInteractiveTextInput.
	// if downloadType == sharedState.TypeVideo {

	// 		WithDefaultValue("https://youtu.be/LXb3EKWsInQ").
	// 		WithDefaultText("URL").
	// 		Show()
	// }

	return makeInput, err
}
