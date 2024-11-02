package sharedState

import (
	"YTDownloaderCli/pkg/_youtube"
	"fmt"
	"strings"

	"github.com/kkdai/youtube/v2"
)

type DownloadType int

const (
	TypeVideo DownloadType = iota
	TypePlaylist
)

// Function to return the name of the download type
func (d DownloadType) String() string {
	return [...]string{"Video", "Playlist"}[d]
}

// DownloadTypeFromString Convert string to enum
func DownloadTypeFromString(s string) (DownloadType, error) {
	switch strings.ToLower(s) {
	case "video":
		return TypeVideo, nil
	case "playlist":
		return TypePlaylist, nil
	default:
		return -1, fmt.Errorf("invalid weekday: %s", s)
	}
}

type SharedState struct {
	DownloadType DownloadType
	URl          string
	//Videos          []videoInfo
	Playlist        []*youtube.PlaylistEntry
	SelectedFormats youtube.Format
	CurrentProgress map[string]progressSate // Map of video IDs to their download progress and other state
	YTclient        _youtube.YTClientModel
}

type progressSate struct {
	Progress float64
	IsDone   bool
}

type videoInfo struct {
	ID          string
	Title       string
	Description string
	Author      string
}
