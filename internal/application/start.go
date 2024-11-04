package application

import (
	"YTDownloaderCli/internal/config"
	"YTDownloaderCli/internal/sharedState"
	"YTDownloaderCli/internal/ui"
	"YTDownloaderCli/internal/utils"
	"YTDownloaderCli/pkg/_youtube"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Start(config *config.Config) {

	defer os.Exit(0)
	F, e := utils.LogToFileWith("log.log", "log")
	if e != nil {
		fmt.Printf("Log File error : %s\n", e)
	}

	state := &sharedState.SharedState{
		YTclient: *_youtube.NewYTClient(),
		Log:      F,
	}
	defer state.Log.Close()

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
}
