package _youtube

import (
	"YTDownloaderCli/internal/utils"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/kkdai/youtube/v2"
)

type YTClientModel struct {
	Client youtube.Client
}

func NewYTClient() *YTClientModel {
	c := &YTClientModel{Client: youtube.Client{ChunkSize: 1024 * 1024 * 2, HTTPClient: &http.Client{Timeout: 2 * time.Minute}}}
	return c

}

func (c *YTClientModel) GetVideoDetail(url string) *youtube.Video {
	videoData, err := c.Client.GetVideo(url)
	utils.UtilError(err)
	return videoData
}

func (c *YTClientModel) GetVideoPlaylistDetail(url string) *youtube.Playlist {
	videoData, err := c.Client.GetPlaylist(url)
	utils.UtilError(err)
	return videoData

}

func (c *YTClientModel) AddPlaylistNumbering(p *youtube.Playlist) *Playlist {
	var playlistEntries []*PlaylistEntry
	for i, video := range p.Videos {
		playlistEntries = append(playlistEntries, &PlaylistEntry{
			PlaylistEntry: video,
			VideoIndex:    i,
		})
	}

	return &Playlist{
		Playlist: p,
		Videos:   playlistEntries,
	}
}

func (c *YTClientModel) GetDownloadStreamWithContext(video *youtube.Video, format *youtube.Format, ctx context.Context) (io.ReadCloser, int64, error) {
	return c.Client.GetStreamContext(ctx, video, format)
}

func (c *YTClientModel) GetDownloadStream(video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
	return c.Client.GetStream(video, format)
}

func (c *YTClientModel) GetVideoFromPlaylistEntry(playlistEntry *PlaylistEntry) (*youtube.Video, error) {
	entry, err := c.Client.VideoFromPlaylistEntry(&youtube.PlaylistEntry{
		ID:         playlistEntry.ID,
		Title:      playlistEntry.Title,
		Author:     playlistEntry.Author,
		Duration:   playlistEntry.Duration,
		Thumbnails: playlistEntry.Thumbnails,
	})
	//utils.UtilError(err)
	return entry, err
}
