package cmd

import (
	utils "YTDownloaderCli/internal/common"
	"YTDownloaderCli/internal/ffmpeg"
	"YTDownloaderCli/internal/yt"

	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() *cobra.Command {

	cmd := &cobra.Command{
		Short: "used to download YT videos",
		Run: func(cmd *cobra.Command, args []string) {
			runApplication()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initViperConfig()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
			// return ffmpegInit()
		},
	}
	cmd.AddCommand(newRun())
	return cmd
}

func runApplication() {
	yt.Start()
}

func initViperConfig() error {

	viper.SetConfigType("toml")

	viper.AddConfigPath(utils.HomeDir())
	viper.AddConfigPath(".")
	viper.AddConfigPath(utils.ExecutableDir())
	viper.SetConfigName("ytdconfig")

	viper.ReadInConfig()
	viper.WatchConfig()

	viper.SetDefault("outputPath", "C:/Users/pp542/Downloads/Video")

	return nil

}

func ffmpegInit() error {

	if !isFFmpegInstalled() {
		ffmpeg.Start()
	}

	return nil
}

func isFFmpegInstalled() bool {

	cmd := exec.Command("ffmpeg", "-version")

	if err := cmd.Run(); err != nil {

		path := viper.GetString("ffmpegPath")

		if path == "" {
			return false
		} else {
			return true
		}
	}
	return true

}
