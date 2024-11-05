package FFMpeg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Chapter represents a video chapter
type Chapter struct {
	StartTime int64  // Start time in milliseconds
	EndTime   int64  // End time in milliseconds (will be calculated based on the next chapter's start time)
	Title     string // Chapter title
}

// parseTime converts HH:MM:SS or MM:SS time strings to milliseconds
func parseTime(timeStr string) (int64, error) {
	parts := strings.Split(timeStr, ":")
	var hours, minutes, seconds int

	if len(parts) == 2 {
		// MM:SS format
		minutes, _ = strconv.Atoi(parts[0])
		seconds, _ = strconv.Atoi(parts[1])
	} else if len(parts) == 3 {
		// HH:MM:SS format
		hours, _ = strconv.Atoi(parts[0])
		minutes, _ = strconv.Atoi(parts[1])
		seconds, _ = strconv.Atoi(parts[2])
	} else {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	// Calculate total milliseconds
	totalMilliseconds := int64((hours*3600 + minutes*60 + seconds) * 1000)
	return totalMilliseconds, nil
}

// ExtractChapters parses the input text and returns a slice of Chapter structs
func ExtractChapters(text string, videoDuration int64) ([]Chapter, error) {
	// Updated regex pattern to handle (HH:MM:SS) or (MM:SS) formats and optional leading or trailing spaces
	pattern := regexp.MustCompile(`\((\d{1,2}:\d{2}(?::\d{2})?)\)\s*(.*)`)
	//pattern := regexp.MustCompile((?m)(?P<time>\d{2}:\d{2}:\d{2}|\d{2}:\d{2})\)?\s(-)?(\s)?(?P<chapterTitle>.*))
	matches := pattern.FindAllStringSubmatch(text, -1)

	var chapters []Chapter
	for _, match := range matches {
		timeStr := match[1] // Captured timestamp
		title := match[2]   // Captured chapter title

		startTime, err := parseTime(timeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %v", err)
		}

		chapters = append(chapters, Chapter{
			StartTime: startTime,
			Title:     title,
		})
	}

	// Set end times for each chapter (next chapter's start time or video duration for the last chapter)
	for i := 0; i < len(chapters)-1; i++ {
		chapters[i].EndTime = chapters[i+1].StartTime
	}
	if len(chapters) > 0 {
		chapters[len(chapters)-1].EndTime = videoDuration
	}

	return chapters, nil
}

// FormatForFFmpeg converts chapters to FFmpeg-compatible metadata format
func FormatForFFmpeg(chapters []Chapter) string {
	var builder strings.Builder

	for _, chapter := range chapters {
		fmt.Fprintf(&builder, "[CHAPTER]\n")
		fmt.Fprintf(&builder, "TIMEBASE=1/1000\n")
		fmt.Fprintf(&builder, "START=%d\n", chapter.StartTime)
		fmt.Fprintf(&builder, "END=%d\n", chapter.EndTime)
		fmt.Fprintf(&builder, "title=%s\n", chapter.Title)
	}

	return builder.String()
}
