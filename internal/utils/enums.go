package utils

type MediaType int

const (
	Unknown MediaType = iota
	Video
	Audio
)

// String method for MediaType to return string representation
func (m MediaType) String() string {
	return [...]string{"Unknown", "Video", "Audio"}[m]
}

// GetMediaType returns the MediaType based on the input string
func GetMediaType(input string) MediaType {
	switch input {
	case "video":
		return Video
	case "audio":
		return Audio
	default:
		return Unknown
	}
}
