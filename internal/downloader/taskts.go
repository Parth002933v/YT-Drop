package downloader

import (
	"YTDownloaderCli/internal/utils"
	"YTDownloaderCli/pkg/_youtube"
	"YTDownloaderCli/pkg/worker"
	"context"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"os"
	"time"
)

type DownloadTask struct {
	id            int
	playlistEntry *youtube.PlaylistEntry
	format        *youtube.Format
	bar           []*mpb.Bar
	YTClient      _youtube.YTClientModel
	p             *mpb.Progress
}

func (t *DownloadTask) Process(ctx context.Context) {
	video := t.YTClient.GetVideoFromPlaylistEntry(t.playlistEntry)
	//utils.GetfprmatInFile(video.Formats)
	if err := finalizeFormat(video.Formats, t.format); err != nil {
		fmt.Printf("error to determine prefered format : %v\n", err)
		return
	}
	mediaType, _, erro := utils.GetVTypeAndCodecFromMimType(t.format.MimeType)
	if erro != nil {
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
		task := fmt.Sprintf("#%02d:", t.id)
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
		videoPath := download(ctx, video, t.format.QualityLabel, "mp4", t.YTClient, t.format, queue[0])

		t.format = utils.GetMaxAudioQuality(video.Formats)
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
		audioPath := download(ctx, video, t.format.AudioQuality, "m4a", t.YTClient, t.format, queue[1])

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

		thumbnailPath, _ := downloadThumbnail(video.Thumbnails[len(video.Thumbnails)-1].URL, video.Title)
		defer os.RemoveAll(videoPath)
		defer os.RemoveAll(audioPath)
		defer os.RemoveAll(thumbnailPath)

		err := mergeVideoAudioThumbnailChapters(videoPath, audioPath, thumbnailPath, video.Title)
		if err != nil {
			queue[2].Abort(true)
			return
		}
		time.Sleep(time.Second * 3)
		queue[2].SetTotal(100, true)
		return
	} else {
		return
	}
	return
}

func Start(playlist []*youtube.PlaylistEntry, format *youtube.Format, client _youtube.YTClientModel) {
	var tasks []worker.Task

	p := mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())

	for i := 0; i < len(playlist); i++ {
		tasks = append(tasks, &DownloadTask{playlistEntry: playlist[i], format: format, YTClient: client, id: i, p: p})
	}

	pool := worker.Pool{
		Tasks:         tasks,
		MaxConcurrent: 2,
	}

	pool.Run()

	p.Wait()
}
