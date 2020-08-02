package room

import "strings"

// Video represents a YouTube video a room can play
type Video struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	Requester *User
}

// NewVideo returns a new video
func NewVideo(url string, requester *User) *Video {
	return &Video{
		URL:       url,
		Requester: requester,
	}
}

func (v *Video) ExtractID() string {
	split := strings.Split(v.URL, "v=")
	if len(split) <= 1 {
		return ""
	}

	videoID := split[1]
	ampersandIdx := strings.Index(videoID, "&")
	if ampersandIdx != -1 {
		videoID = videoID[:ampersandIdx]
	}

	return videoID
}

// SaveVideoRequest is a request to save a video
type SaveVideoRequest struct {
	Video *Video
	Room  *Room
}
