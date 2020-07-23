package room

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
