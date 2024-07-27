package ytclient

import (
	"YTDownloaderCli/internal/utils"
	"io"

	"github.com/kkdai/youtube/v2"
)

type YTClientModel struct {
	Client youtube.Client
}

func NewYTClient() *YTClientModel {
	c := &YTClientModel{Client: youtube.Client{ChunkSize: 1024 * 1024 * 5}}
	return c
}

func (c *YTClientModel) GetVideoDetail(url string) *youtube.Video {

	videoData, err := c.Client.GetVideo(url)
	utils.PError(err)

	return videoData
}

func (c *YTClientModel) GetVideoPlaylistDetail(url string) *youtube.Playlist {

	videoData, err := c.Client.GetPlaylist(url)

	utils.PError(err)

	return videoData

}

func (c *YTClientModel) GetDownloadStream(video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
	return c.Client.GetStream(video, format)
}
