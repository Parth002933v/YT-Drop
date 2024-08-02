package ffmpeg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	utils "YTDownloaderCli/internal/common"
	"YTDownloaderCli/internal/yt"

	tea "github.com/charmbracelet/bubbletea"
)

func USelectionStage(m *ffmpegModel, msg tea.Msg, cmd *tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case tea.KeyRight.String(), tea.KeyLeft.String():
			m.selectionOption.wantToDownload = !m.selectionOption.wantToDownload

		case tea.KeyEnter.String():
			m.questions.defaultVal = utils.HomeDir()
			if m.selectionOption.wantToDownload {
				downloadFFMpegUpdate(m, msg, cmd)
			}
			m.stage = 2
		}
	}
}

func UPostSelection(m *ffmpegModel, msg tea.Msg, cmd *tea.Cmd) {
	if !m.selectionOption.wantToDownload {
		specifyFFMpegPathUpdate(m, msg, cmd)
	} else {
	}
}

func specifyFFMpegPathUpdate(m *ffmpegModel, msg tea.Msg, cmd *tea.Cmd) {
	current := m.questions
	if !m.selectionOption.wantToDownload {
		switch msg := msg.(type) {
		case tea.KeyMsg:

			switch msg.String() {
			case tea.KeyTab.String():
				m.questions.input.SetValue(m.questions.defaultVal)

			case tea.KeyCtrlC.String():
				*cmd = tea.Quit

			case tea.KeyEnter.String():
				// m.questions.input.Blur()
				// m.stage = 3
				// installFFmpegManuallyWindows(m)
			}
		}
		current.input, *cmd = current.input.Update(msg)

	}
}
func downloadFFMpegUpdate(m *ffmpegModel, msg tea.Msg, cmd *tea.Cmd) {
	//todo
	current := m.questions
	if !m.selectionOption.wantToDownload {
		switch msg := msg.(type) {
		case tea.KeyMsg:

			switch msg.String() {
			case tea.KeyTab.String():
				m.questions.input.SetValue(m.questions.defaultVal)

			case tea.KeyCtrlC.String():
				*cmd = tea.Quit

			case tea.KeyEnter.String():
				m.questions.input.Blur()
				// m.stage = 3
				installFFmpegManuallyWindows(m)
			}
		}
		current.input, *cmd = current.input.Update(msg)

	}
}

//*==================================================================================

func installFFmpegManuallyWindows(m *ffmpegModel) {
	m.downloader.isProcessing = true

	url := "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
	resp, err := http.Get(url)
	utils.UtilError(err)
	defer resp.Body.Close()

	outFile := fmt.Sprintf("%s\\%s", utils.HomeDir(), "ffmpeg.zip")
	fmt.Print(outFile)

	out, err := os.Create(outFile)
	utils.UtilError(err)
	defer out.Close()

	pw := &yt.ProgressWriter{
		Writer:     out,
		TotalBytes: resp.ContentLength,
		ProgressDisplay: func(progress float64) {
			m.downloader.downloadPrecentage = progress
		},
	}

	_, err = io.Copy(pw, resp.Body)
	utils.UtilError(err)

	m.downloader.isProcessing = false

	m.downloader.isProcessing = true
	err = unzip(outFile, "./")
	if err != nil {
		utils.UtilError(fmt.Errorf("error extracting ffmpeg: %v", err))
		return
	}

	m.downloader.isProcessing = false
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "bin/ffmpeg.exe" {
			fpath := filepath.Join(dest, f.Name)
			if f.FileInfo().IsDir() {
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}
