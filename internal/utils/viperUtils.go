package utils

import (
	"YTDownloaderCli/internal/config"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateDefaultViperConfigFile(confiigFilePath string) {
	viper.SetDefault("app-info.app_name", "YTDownloaderCli")
	viper.SetDefault("download-preference.video_pref", "avc1")
	viper.SetDefault("download-preference.audio_pref", "mp4a")

	if err := viper.WriteConfigAs(confiigFilePath); err != nil {
		cobra.CheckErr(fmt.Sprintf("Error creating config file: %v \n", err))
	}
}

func UnmarshalViperToConfig(cfg *config.Config) {
	if err := viper.Unmarshal(&cfg); err != nil {
		cobra.CheckErr(fmt.Sprintf("Error in Unmarshaling config: %v", err))
	}
}
