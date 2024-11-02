package application

import (
	"YTDownloaderCli/internal/downloader"
	"YTDownloaderCli/internal/sharedState"
	"YTDownloaderCli/internal/ui"
	"YTDownloaderCli/internal/ui/theme"
	"YTDownloaderCli/internal/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func handlePlaylist(state *sharedState.SharedState) {

	spinner := ui.Spinner()
	res := state.YTclient.GetVideoPlaylistDetail(state.URl)
	spinner.Success()

	// print playlist detail
	fmt.Println(playlistDetailPrint(res))

	// video selection
	selectedVideos := ui.MultiVideoSelection(res.Videos)
	state.Playlist = selectedVideos

	//format lists
	formats := utils.GetFormats()
	formats.Sort()
	//utils.FilterFormatsByMineType(formats, "vp9", "opus")
	//utils.GetfprmatInFile(*formats)
	//formats selection
	selectedFormat := ui.FormatSelection(*formats)
	state.SelectedFormats = selectedFormat

	downloader.Start(state.Playlist, &state.SelectedFormats, state.YTclient)
	fmt.Println("end of downloader.Start()")
}

func playlistDetailPrint(video *youtube.Playlist) string {
	var metadataBuilder strings.Builder
	metadataBuilder.WriteString(fmt.Sprintf("Author   : %s\n", video.Author))
	metadataBuilder.WriteString(fmt.Sprintf("Title    : %s\n", video.Title))
	metadataBuilder.WriteString(fmt.Sprintf("VideoCount    : %d", len(video.Videos)))

	styledText := theme.BoxTheme().Render(metadataBuilder.String())

	return styledText
}

func GetCodeFromMimType(str string) (codec string) {
	// Regular expression to capture the codec identifier (e.g., avc1, av01, vp9)
	re := regexp.MustCompile(`codecs="([a-z0-9]+)`)

	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		return match[1]
	} else {
		return "Unknown Codec"
	}
}
