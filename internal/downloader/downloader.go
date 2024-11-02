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

func downloadThumbnail(url string, name string) (string, error) {
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

func finalizeFormat(videoFormats youtube.FormatList, selectedFormat *youtube.Format) {
	staticFormats := utils.GetFormats()
	staticFormats.Sort()

	//utils.FilterFormatsByMineType(staticFormats, "vp9", "opus")

	// Loop until the desired format is found in video.Formats
	for {
		//selected format in available in fetched video
		format, _, ok := isFormatAvailable(videoFormats, *selectedFormat)
		if ok {
			selectedFormat = &format
			// Desired format found, break out of the loop
			break
		}

		// If the format is not found, check in static formats
		_, i, foundInStatic := isFormatAvailable(*staticFormats, *selectedFormat)
		if !foundInStatic {
			return
		}

		// Check if we can fall back to a lower format
		if i+1 < len(*staticFormats) {

			fmt.Printf("Format not found, falling back to lower format with ItagNo: %v\n", (*staticFormats)[i+1].ItagNo)
			selectedFormat.ItagNo = (*staticFormats)[i+1].ItagNo // Fallback to previous format
		} else {
			// No previous format to fall back to, exit with error
			return
		}
	}
}
