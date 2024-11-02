package downloader

import (
	"YTDownloaderCli/internal/utils"
	"YTDownloaderCli/pkg/_youtube"
	"context"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"github.com/vbauerster/mpb/v8"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func download(ctx context.Context, video *youtube.Video, suffix string, extension string, client _youtube.YTClientModel, format *youtube.Format, bar *mpb.Bar) (filePath string) {

	stream, n, err := client.GetDownloadStreamWithContext(video, format, ctx)
	if err != nil {
		fmt.Printf("Failed to get download stream: %v\n", err)
		return
	}
	defer stream.Close()

	bar.SetTotal(n, false)
	defer bar.SetTotal(n, true)

	// Monitor context for cancellation and abort the bar if needed
	go func() {
		for {
			select {
			case <-ctx.Done():
				bar.Abort(true)
				return
			default:
				time.Sleep(100 * time.Millisecond)
				// Continue monitoring
			}
		}
	}()

	reader := bar.ProxyReader(stream)
	defer reader.Close()

	// Create file for saving the download
	_file, err := MakeFile(video.Title, suffix, extension)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer _file.Close()

	_, err = io.Copy(_file, reader)
	if err != nil {
		fmt.Printf("Failed to copy stream to file: %v\n", err)
		return
	}

	return _file.Name()
}

func downloadThumbnail(url string, name string) (thumbnailPath string, erro error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	file, erro := MakeFile(downloader.SanitizeFilename(name), "thumbnail", "jpg")
	if erro != nil {
		return "", erro
	}
	defer file.Close()

	io.Copy(file, res.Body)
	return file.Name(), nil

}

func MakeFile(title, suffix string, extension string) (file *os.File, erro error) {
	fileName := fmt.Sprintf("%s_%s.%v", downloader.SanitizeFilename(title), suffix, extension)
	return os.Create(downloader.SanitizeFilename(fileName))
}

func isFormatAvailable(formatList youtube.FormatList, targetedFormat youtube.Format) (format youtube.Format, i int, ok bool) {
	for i, f := range formatList {
		if f.ItagNo == targetedFormat.ItagNo {
			return f, i, true
		}
	}
	return youtube.Format{}, -1, false
}

func finalizeFormat(videoFormats youtube.FormatList, selectedFormat *youtube.Format) error {
	staticFormats := utils.GetFormats()
	staticFormats.Sort()

	// Loop until the desired format is found in videoFormats
	for {
		// Check if the selected format is available in the fetched video formats
		format, _, available := isFormatAvailable(videoFormats, *selectedFormat)
		if available {
			*selectedFormat = format // Update the actual content of selectedFormat
			return nil               // Exit function as format has been found
		}

		// If the format is not found, check in static formats
		_, i, foundInStatic := isFormatAvailable(*staticFormats, *selectedFormat)
		if !foundInStatic {
			return fmt.Errorf("format not available in static formats")
		}
		//if i+1 < len(*staticFormats) {
		// Attempt to fall back to a lower format (if sorted in ascending quality)
		if i > 0 { // Check if there’s a lower quality option available
			fmt.Printf("Format not found, falling back to lower format with ItagNo: %v\n", (*staticFormats)[i-1].ItagNo)
			*selectedFormat = (*staticFormats)[i-1] // Fallback to lower format
		} else {
			// No lower format available, exit with error
			return fmt.Errorf("no suitable format found")
		}
	}
}

// mergeVideoAudioThumbnailChapters merges a video and audio file using FFmpeg
func mergeVideoAudioThumbnailChapters(videoPath, audioPath, thumbnailPath, outputName string) error {

	ffmpegPath2, err := utils.ExtractFFmpeg()
	utils.UtilError(err)
	defer os.RemoveAll(filepath.Dir(ffmpegPath2)) // Clean up temporary directory

	outputFileName := fmt.Sprintf("%v.mp4", downloader.SanitizeFilename(outputName))
	// // Prepare FFmpeg command arguments
	// args := []string{
	// 	"-y",
	// 	"-i", videoPath,
	// 	"-i", thumbnailPath,
	// 	"-i", audioPath,
	// 	"-map", "0:v",
	// 	"-map", "1",
	// 	"-map", "2:a",
	// 	"-c:v", "copy",
	// 	"-c:a", "copy",
	// 	"-c:v:1", "png",
	// 	"-disposition:v:1", "attached_pic",
	// 	outputFileName,
	// }

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

	args = append(args, outputFileName)

	// Create and execute the command
	cmd := exec.Command(ffmpegPath2, args...)

	// Capture output
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg command failed: %v", err)
	}

	// F.WriteString("Video processing completed successfully.\n")
	// F.WriteString(fmt.Sprintf("FFmpeg output: %s\n", output))
	return nil
}
