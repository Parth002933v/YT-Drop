package ui

import "github.com/charmbracelet/lipgloss"

func Logo() string {
	const title = `
	__   _______   ____
	\ \ / /_   _| |  _ \ _ __ ___  _ __
	 \ V /  | |   | | | |  __/ _ \|  _ \
	  | |   | |   | |_| | | | (_) | |_) |
	  |_|   |_|   |____/|_|  \___/|  __/
								  |_|
	`
	titleStyle := lipgloss.NewStyle().
		Bold(true).Padding(0).Margin(0).
		Foreground(lipgloss.Color("#FF6347")) // Tomato color
	return titleStyle.Render(title)
}

func SubTitle() string {
	const subtitle = "Your ultimate YouTube video downloader!"
	subtitleStyle := lipgloss.NewStyle().
		Italic(true).Padding(0).Margin(0).
		Foreground(lipgloss.Color("#87CEEB")) // SkyBlue color

	return subtitleStyle.Render(subtitle)
}
