package application

import (
	"YTDownloaderCli/internal/downloader"
	"YTDownloaderCli/internal/sharedState"
	"YTDownloaderCli/internal/ui"
	"YTDownloaderCli/internal/ui/theme"
	"YTDownloaderCli/pkg/_youtube"
	"strings"

	"fmt"

	"github.com/kkdai/youtube/v2"
)

func handleVideo(state *sharedState.SharedState) {

	//fetch video data along with formats
	spinner := ui.Spinner()
	video := state.YTclient.GetVideoDetail(state.URl)
	spinner.Success()

	// print the video info in bix
	fmt.Println(videoDetailPrint(*video))

	video.Formats.Sort()
	//utils.FilterFormatsByMineType(&video.Formats, "vp9", "opus")
	//utils.GetfprmatInFile(video.Formats)

	//format selection
	selectedFormat := ui.FormatSelection(video.Formats)

	state.SelectedFormats = selectedFormat

	downloader.Start(
		[]*_youtube.PlaylistEntry{
			{
				PlaylistEntry: &youtube.PlaylistEntry{
					ID:         video.ID,
					Title:      video.Title,
					Author:     video.Author,
					Duration:   video.Duration,
					Thumbnails: video.Thumbnails,
				},
			},
		},
		&state.SelectedFormats,
		state.YTclient,
		state.Log,
	)

}

func videoDetailPrint(video youtube.Video) string {
	var metadataBuilder strings.Builder
	metadataBuilder.WriteString(fmt.Sprintf("Author   : %s\n", video.Author))
	metadataBuilder.WriteString(fmt.Sprintf("Title    : %s\n", video.Title))
	metadataBuilder.WriteString(fmt.Sprintf("Duration : %s\n", video.Duration))
	metadataBuilder.WriteString(fmt.Sprintf("Views    : %v", video.Views))

	styledText := theme.BoxTheme().Render(metadataBuilder.String())

	return styledText
}
