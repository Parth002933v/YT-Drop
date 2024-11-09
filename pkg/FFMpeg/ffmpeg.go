package FFMpeg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Chapter represents a video chapter
type Chapter struct {
	StartTime int64
	EndTime   int64
	Title     string
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

	// Calculate milliseconds
	totalMilliseconds := int64((hours*3600 + minutes*60 + seconds) * 1000)
	return totalMilliseconds, nil
}

// ExtractChapters parses the input text and returns a slice of Chapter structs
func ExtractChapters(text string, videoDuration int64) ([]Chapter, error) {
	pattern := regexp.MustCompile(`(?m)^\s*[^\w\n]*\(?(\d{1,2}:\d{2}(?::\d{2})?)\)?\s*[^\w\n]*\s*(.+)$`)
	matches := pattern.FindAllStringSubmatch(text, -1)

	var chapters []Chapter
	for _, match := range matches {
		timeStr := match[1]
		title := strings.TrimSpace(match[2])

		startTime, err := parseTime(timeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %v", err)
		}

		chapters = append(chapters, Chapter{
			StartTime: startTime,
			Title:     title,
		})
	}

	// Validate chapter order and assign end times
	for i := 0; i < len(chapters)-1; i++ {
		current := &chapters[i]
		next := &chapters[i+1]

		if next.StartTime <= current.StartTime {
			return nil, fmt.Errorf("invalid chapter timing: chapter %d starts at %dms but previous chapter ends at %dms", i+1, next.StartTime, current.StartTime)
		}
		current.EndTime = next.StartTime // Set end time based on next chapter's start time
	}

	// Set end time for the last chapter based on video duration
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
