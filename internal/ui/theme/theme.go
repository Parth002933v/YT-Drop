package theme

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func PrintErrorText(text string) {

	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6347"))

	fmt.Print(errorStyle.Render(text))
}

func BoxTheme() lipgloss.Style {

	metadataBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		PaddingLeft(1).PaddingRight(1)

	return metadataBox
}
