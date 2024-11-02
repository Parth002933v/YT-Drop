package cmd

import (
	"YTDownloaderCli/pkg/worker"
	"context"
	cr "crypto/rand"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

func init() {
	NewRun().AddCommand(textRun())
}
func textRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "short test of how to use it",
		Run: func(cmd *cobra.Command, args []string) {
			main3()
		},
	}
	return cmd
}

func main3() {
	var wg sync.WaitGroup
	client := &youtube.Client{}
	wg.Add(3)
	go down(client, &wg, "Xgs2WCHkmbQ")
	go down(client, &wg, "1PxhTfmEyQ8")
	go down(client, &wg, "BaW_jenozKc")

	wg.Wait()
}

func down(client *youtube.Client, wg *sync.WaitGroup, s string) {
	defer wg.Done()
	videoID := s

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	file, err := os.Create(fmt.Sprintf("%s.mp4", videoID))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
}

type testTask struct {
	vtype int
	p     *mpb.Progress
}

func (t *testTask) Process(ctx context.Context) {
	if t.vtype == 1 {
		bar := t.p.New(0,
			mpb.BarStyle().Rbound("|"),
			mpb.PrependDecorators(
				decor.Name("playlist-id", decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("task", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
				decor.OnCompleteOrOnAbort(
					decor.Counters(
						decor.SizeB1024(0), "% .2f / % .2f"),
					"Done!",
				),
			),
		)

		simulateDownload(bar) // simulate video download
	} else {
		queue := make([]*mpb.Bar, 2)

		queue[0] = t.p.New(0,
			mpb.BarStyle().Rbound("|"),
			mpb.PrependDecorators(
				decor.Name("playlist-id", decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.Name("task", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
				decor.OnCompleteOrOnAbort(
					decor.Counters(
						decor.SizeB1024(0), "% .2f / % .2f"),
					"Done!",
				),
			),
		)

		simulateDownload(queue[0]) // simulate video download
		queue[1] = t.p.New(0,
			mpb.BarStyle().Rbound("|"),
			mpb.BarQueueAfter(queue[0]),
			mpb.PrependDecorators(
				decor.Name("task", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
				decor.Name("Playlist-id", decor.WCSyncWidth, decor.WCSyncSpaceR),
				decor.OnCompleteOrOnAbort(
					decor.Counters(
						decor.SizeB1024(0), "% .2f / % .2f"),
					"Done!",
				),
			),
		)

		// Download audio if video format is chosen
		simulateDownload(queue[1]) // simulate audio download	j
	}

}

// simulateDownload simulates a download task by updating progress.
func simulateDownload(bar *mpb.Bar) {
	var total int64 = 64 * 1024 * 1024

	bar.SetTotal(total, false)
	defer bar.SetTotal(total, true)
	r, w := io.Pipe()

	go func() {

		for i := 0; i < 1024; i++ {
			_, _ = io.Copy(w, io.LimitReader(cr.Reader, 64*1024))
			sleepDuration := time.Duration(rand.Intn(40)+10) * time.Millisecond
			time.Sleep(sleepDuration)

		}
		w.Close()
	}()
	proxyReader := bar.ProxyReader(r)
	defer proxyReader.Close()

	// copy from proxyReader, ignoring errors
	_, _ = io.Copy(io.Discard, proxyReader)
	//total := 100
	//
	//bar.SetTotal(100, false)
	//for i := 0; i < total; i++ {
	//	bar.Increment()
	//	time.Sleep(time.Duration(rand.Intn(80)) * time.Millisecond) // simulate variable download speed
	//}
	//bar.SetTotal(100, true)
}
func main2() {
	var tasks []worker.Task

	p := mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())

	for i := 0; i < 3; i++ {
		tasks = append(
			tasks,
			&testTask{vtype: i, p: p},
		)
	}
	pool := worker.Pool{
		Tasks:         tasks,
		MaxConcurrent: 2,
	}

	pool.Run()

	p.Wait()

}
