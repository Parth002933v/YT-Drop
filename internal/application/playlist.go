package application

import (
	"YTDownloaderCli/internal/downloader"
	"YTDownloaderCli/internal/sharedState"
	"YTDownloaderCli/internal/ui"
	"YTDownloaderCli/internal/ui/theme"
	"YTDownloaderCli/internal/utils"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"strings"
)

func handlePlaylist(state *sharedState.SharedState) {

	spinner := ui.Spinner()
	playlist := state.YTclient.GetVideoPlaylistDetail(state.URl)
	spinner.Success()

	// print playlist detail
	fmt.Println(playlistDetailPrint(playlist))

	newPlaylist := state.YTclient.AddPlaylistNumbering(playlist)

	// video selection
	selectedVideos := ui.MultiVideoSelection(newPlaylist.Videos)
	state.Playlist = selectedVideos

	//format lists
	formats := utils.GetFormats()
	formats.Sort()
	//utils.FilterFormatsByMineType(formats, "vp9", "opus")
	//utils.GetfprmatInFile(*formats)
	//formats selection
	selectedFormat := ui.FormatSelection(*formats)
	state.SelectedFormats = selectedFormat

	downloader.Start(state.Playlist, &state.SelectedFormats, state.YTclient, nil)
}

func playlistDetailPrint(video *youtube.Playlist) string {
	var metadataBuilder strings.Builder
	metadataBuilder.WriteString(fmt.Sprintf("Author   : %s\n", video.Author))
	metadataBuilder.WriteString(fmt.Sprintf("Title    : %s\n", video.Title))
	metadataBuilder.WriteString(fmt.Sprintf("VideoCount    : %d", len(video.Videos)))

	styledText := theme.BoxTheme().Render(metadataBuilder.String())

	return styledText
}
