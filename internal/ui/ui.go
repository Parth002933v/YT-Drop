package ui

import (
	"YTDownloaderCli/internal/utils"

	"github.com/pterm/pterm"
)

func Spinner() *pterm.SpinnerPrinter {

	spinner1, err := pterm.DefaultSpinner.
		WithRemoveWhenDone(true).
		WithShowTimer(false).
		WithText("Loading...").
		Start()

	utils.UtilError(err)
	return spinner1
}
