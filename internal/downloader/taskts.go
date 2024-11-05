package downloader

import (
	"YTDownloaderCli/internal/utils"
	"YTDownloaderCli/pkg/_youtube"
	"YTDownloaderCli/pkg/worker"
	"context"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"os"
	"runtime"
)

type DownloadTask struct {
	id            *int
	playlistEntry *_youtube.PlaylistEntry
	format        *youtube.Format
	bar           []*mpb.Bar
	YTClient      _youtube.YTClientModel
	p             *mpb.Progress
	log           *os.File
}

func (t *DownloadTask) Process(ctx context.Context) {
	var task string
	if t.id != nil {
		task = fmt.Sprintf("#%02d:", *t.id)
	} else {
		task = "#00:"
	}

	video, erro := t.YTClient.GetVideoFromPlaylistEntry(t.playlistEntry)
	if erro != nil {
		t.log.WriteString(fmt.Sprintf("error in video retrival : %s\n", erro))
	}

	//utils.GetfprmatInFile(video.Formats)
	if err := finalizeFormat(video.Formats, t.format); err != nil {
		t.log.WriteString(fmt.Sprintf("task: %02d error to determine prefered format : %v\n", task, err))
		fmt.Printf("task: %02d error to determine prefered format : %v\n", task, err)
		return
	}
	mediaType, _, erro := utils.GetVTypeAndCodecFromMimType(t.format.MimeType)
	if erro != nil {
		t.log.WriteString(fmt.Sprintf("error to determine mimetype : %v\n", erro))
		return
	}

	if mediaType == utils.Audio {
		task := fmt.Sprintf("#%02d:", t.id)
		bar := t.p.New(0,
			mpb.BarStyle().Rbound("|"),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(task, decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("0/1", decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name(utils.TruncateWithEllipsisString(t.playlistEntry.Title, 30), decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
			mpb.AppendDecorators(
				decor.OnCompleteOrOnAbort(
					decor.Counters(
						decor.SizeB1024(0), "% .2f / % .2f"),
					"Done!",
				),
			),
		)
		_ = download(ctx, video, t.format.AudioQuality, "m4a", t.YTClient, t.format, bar)

	} else if mediaType == utils.Video {
		queue := make([]*mpb.Bar, 3)

		queue[0] = t.p.AddBar(0,
			mpb.PrependDecorators(
				decor.Name(task, decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("1/3", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
			mpb.AppendDecorators(
				decor.Counters(
					decor.SizeB1024(0), "% .2f / % .2f"),
			),
		)

		queue[1] = t.p.AddBar(0,
			mpb.BarQueueAfter(queue[0]),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(task, decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("2/3", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
			mpb.AppendDecorators(
				decor.Counters(
					decor.SizeB1024(0), "% .2f / % .2f"),
			),
		)

		queue[2] = t.p.New(0,
			mpb.SpinnerStyle("∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"),
			mpb.BarQueueAfter(queue[1]),
			mpb.BarFillerClearOnComplete(),
			mpb.BarFillerOnComplete("Done!"),
			mpb.PrependDecorators(
				decor.Name(task, decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("3/3", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
		)

		videoPath := download(ctx, video, t.format.QualityLabel, "mp4", t.YTClient, t.format, queue[0])
		t.format = utils.GetMaxAudioQuality(video.Formats)
		audioPath := download(ctx, video, t.format.AudioQuality, "m4a", t.YTClient, t.format, queue[1])
		thumbnailPath, _ := downloadThumbnail(video.Thumbnails[len(video.Thumbnails)-1].URL, video.Title)
		chaptersPath, _ := downloadChapters(video.Description, downloader.SanitizeFilename(video.Title), video.Duration)

		defer os.RemoveAll(chaptersPath)
		defer os.RemoveAll(videoPath)
		defer os.RemoveAll(audioPath)
		defer os.RemoveAll(thumbnailPath)

		queue[2].SetTotal(0, false)

		outputFileName := ""
		if t.id != nil {
			outputFileName = fmt.Sprintf("%02d %v", t.playlistEntry.VideoIndex+1, downloader.SanitizeFilename(video.Title))
		} else {
			outputFileName = downloader.SanitizeFilename(video.Title)
		}
		err := mergeVideoAudioThumbnailChapters(videoPath, audioPath, thumbnailPath, chaptersPath, outputFileName, t.log)
		if err != nil {
			t.log.WriteString(fmt.Sprintf("FFmpeg command failed: %v\n", err))
			fmt.Printf("processing failed for task : %v \n", task)
			queue[2].Abort(true)
			return
		}
		queue[2].SetTotal(100, true)
		return
	} else {
		return
	}
	return
}

func Start(playlist []*_youtube.PlaylistEntry, format *youtube.Format, client _youtube.YTClientModel, log *os.File) {
	var tasks []worker.Task

	p := mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())

	for i := 0; i < len(playlist); i++ {
		if len(playlist) == 1 {
			tasks = append(tasks, &DownloadTask{playlistEntry: playlist[i], format: format, YTClient: client, id: nil, p: p, log: log})
		} else {
			tasks = append(tasks, &DownloadTask{playlistEntry: playlist[i], format: format, YTClient: client, id: &i, p: p, log: log})
		}
	}

	cpuCount := runtime.NumCPU() / 2
	if cpuCount < 2 {
		cpuCount = 1
	}

	fmt.Println(cpuCount)
	pool := worker.Pool{
		Tasks:         tasks,
		MaxConcurrent: cpuCount,
	}

	pool.Run()

	p.Wait()
}
