package application

import (
	utils "YTDownloaderCli/internal/common"
	"YTDownloaderCli/internal/config"
	"YTDownloaderCli/internal/sharedState"
	"YTDownloaderCli/internal/ui"
	"YTDownloaderCli/pkg/_youtube"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Start(config *config.Config) {

	defer os.Exit(0)
	state := &sharedState.SharedState{
		YTclient: *_youtube.NewYTClient(),
	}

	fmt.Println(ui.Logo())
	fmt.Println(ui.SubTitle())

	// select the url type video or playlist
	fmt.Println()
	selectedOption, err := ui.TypeSelection()
	cobra.CheckErr(err)
	state.DownloadType = selectedOption

	//provide url
	URL, err := ui.UrlInput(state.DownloadType)
	cobra.CheckErr(err)
	state.URl = URL

	switch state.DownloadType {
	case sharedState.TypeVideo:
		handleVideo(state)
	case sharedState.TypePlaylist:
		handlePlaylist(state)
	default:
		utils.UtilError(errors.New("invalid url"))
	}
	//fmt.Println("end of start.go")
}
