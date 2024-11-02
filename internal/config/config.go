package config

type Config struct {
	AppInfo            appInfo            `mapstructure:"app-info"`
	DownloadPreference downloadPreference `mapstructure:"download-preference"`
}

type appInfo struct {
	AppName string `mapstructure:"app_name"`
	Verison string `mapstructure:"verison"`
}

type downloadPreference struct {
	VideoPreferce string `mapstructure:"video_pref"`
	AudioPreferce string `mapstructure:"audio_pref"`
}
