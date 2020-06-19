package room

import "net/url"

// Video represents a YouTube video a room can play
type Video struct {
	URL   *url.URL
	Title string
}
