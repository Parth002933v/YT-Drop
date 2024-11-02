package FFMpeg

import (
	"YTDownloaderCli/internal/utils"
	"fmt"
	"github.com/kkdai/youtube/v2/downloader"
	"os"
	"os/exec"
)

func MergeAudioVideo(videoPath, thumbnailPath, audioPath, outPutFileName string) {

	ffmpegPath, err := utils.ExtractFFmpeg()
	utils.UtilError(err)
	defer os.RemoveAll(ffmpegPath)

	outPutFile := fmt.Sprintf("%v.mp4", downloader.SanitizeFilename(outPutFileName))

	args := []string{
		"-y",
		"-i", videoPath,
		"-i", thumbnailPath,
		"-i", audioPath,
	}

	// Conditionally add chaptersPath if it's not nil
	// if chaptersPath != nil && *chaptersPath != "" {
	// 	args = append(args, "-i", *chaptersPath)
	// }

	// Add the remaining arguments
	args = append(args,
		"-map", "0:v",
		"-map", "1",
		"-map", "2:a",
		"-c:v", "copy",
		"-c:a", "copy",
		"-c:v:1", "png",
		"-disposition:v:1", "attached_pic",
	)

	// If chaptersPath was added, map metadata
	// if chaptersPath != nil && *chaptersPath != "" {
	// 	args = append(args, "-map_metadata", fmt.Sprintf("%d", len(args)-1))
	// }

	args = append(args, outPutFile)

	// Create and execute the command
	cmd := exec.Command(ffmpegPath, args...)

	_, err = cmd.CombinedOutput()
}
