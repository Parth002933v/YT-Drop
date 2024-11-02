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
	fmt.Printf("%+v\n", t.format)
	fmt.Printf("%+v\n", t.format)
	finalizeFormat(video.Formats, t.format)
	fmt.Printf("%+v\n", t.format)
	fmt.Printf("%+v\n", t.format)
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
		_ = download(ctx, video, t.format.QualityLabel, "mp4", t.YTClient, t.format, queue[0])

		t.format = utils.GetMaxAudioQuality(video.Formats)
		queue[1] = t.p.AddBar(0,
			mpb.BarQueueAfter(queue[0]),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(task, decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("2/3", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
			mpb.AppendDecorators(
				decor.OnCompleteOrOnAbort(
					decor.Counters(
						decor.SizeB1024(0), "% .2f / % .2f"),
					"Done!",
				),
			),
		)
		_ = download(ctx, video, t.format.AudioQuality, "m4a", t.YTClient, t.format, queue[1])

		queue[2] = t.p.New(0,
			mpb.SpinnerStyle("∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"),
			mpb.BarQueueAfter(queue[1]),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(task, decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("3/3", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
		)
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
		tasks = append(tasks, &DownloadTask{playlistEntry: playlist[i], format: format, YTClient: *_youtube.NewYTClient(), id: i, p: p})
	}

	pool := worker.Pool{
		Tasks:         tasks,
		MaxConcurrent: 2,
	}

	pool.Run()

	p.Wait()
}
