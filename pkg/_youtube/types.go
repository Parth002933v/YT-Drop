package _youtube

import (
	"github.com/kkdai/youtube/v2"
)

//	type TYouTubePlaylistEntry struct {
//		*youtube.PlaylistEntry
//		VideoIndex int
//	}
type Playlist struct {
	*youtube.Playlist // original playlist type
	Videos            []*PlaylistEntry
}

type PlaylistEntry struct {
	*youtube.PlaylistEntry // original playlistEntry type
	VideoIndex             int
}
