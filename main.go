package main

import (
	"YTDownloaderCli/cmd"
	"context"
	"os"
	"os/signal"
)



func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	if err := cmd.New().ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
