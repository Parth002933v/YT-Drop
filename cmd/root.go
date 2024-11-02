package cmd

import (
	// "YTDownloaderCli/internal/config"

	"YTDownloaderCli/internal/application"
	"YTDownloaderCli/internal/config"
	"YTDownloaderCli/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	// "github.com/fsnotify/fsnotify"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewRun() *cobra.Command {
	appConfig := config.Config{}
	runCmd := &cobra.Command{
		Use:   "YTDownloaderCli",
		Short: "Start application",
		Long:  "this application is used to download _youtube videos and playlist though cli",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initViperConfig(cmd, &appConfig)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runApplication(&appConfig)
			//main3()
			//main2()
		},
	}
	return runCmd
}

func initViperConfig(cmd *cobra.Command, cfg *config.Config) error {
	// gets the home directory path where the config file should be located
	configDir := utils.HomeDir()

	//create the config direcory if not exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			cobra.CheckErr(fmt.Sprintf("Error creating config directory: %v\n", err))
		}
	}

	//set config file path in viper
	configFilePath := filepath.Join(configDir, "config.yaml")
	viper.SetConfigFile(configFilePath)

	// check if the config file is exist or not. if not create it with default values
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		utils.CreateDefaultViperConfigFile(configFilePath)
		fmt.Println("Configuration file created with default settings at:", configFilePath)
	}

	//read in config
	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(fmt.Sprintf("Error reading config file: %v", err))
	}

	utils.UnmarshalViperToConfig(cfg)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Print("OnConfigChange")
		utils.UnmarshalViperToConfig(cfg)
	})

	// bindFlags(cmd.Flags(), "", reflect.ValueOf(config.Config{}))
	return nil
}

func runApplication(cfg *config.Config) {
	app := fx.New(
		fx.Provide(func() *config.Config {
			return cfg
		}),

		fx.NopLogger,
		fx.Invoke(application.Start),
	)

	app.Run()

}

func bindFlags(flags *pflag.FlagSet, prefix string, v reflect.Value) {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := range t.NumField() {
		field := t.Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			bindFlags(flags, fmt.Sprintf("%s.%s", prefix, strings.ToLower(field.Name)), v.Field(i))
		default:
			newPrefix := prefix[1:]
			newName := modifyFlag(field.Name)
			configName := fmt.Sprintf("%s.%s", newPrefix, newName)
			flag := flags.Lookup(fmt.Sprintf("%s-%s", strings.ReplaceAll(newPrefix, ".", "-"), newName))
			if !flag.Changed && viper.IsSet(configName) {
				confVal := viper.Get(configName)
				if field.Type.Kind() == reflect.Slice {
					sliceValue, ok := confVal.([]interface{})
					if ok {
						for _, v := range sliceValue {
							flag.Value.Set(fmt.Sprintf("%v", v))
						}
					}
				} else {
					flags.Set(flag.Name, fmt.Sprintf("%v", confVal))
				}
			}
		}
	}
}
func modifyFlag(s string) string {
	var result []rune

	for i, c := range s {
		if i > 0 && unicode.IsUpper(c) {
			result = append(result, '-')
		}
		result = append(result, unicode.ToLower(c))
	}

	return string(result)
}
