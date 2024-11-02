package main

import (
	// "YTDownloaderCli/cmd"
	// "YTDownloaderCli/internal/yt"
	"YTDownloaderCli/cmd"
	"context"
	"os"
	"os/signal"
	// "context"
	// "os"
	// "os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cmd.NewRun().ExecuteContext(ctx); err != nil {
		panic(err)
	}
	// yt.Start()
}
