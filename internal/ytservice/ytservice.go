package ytservice

import (
	"context"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var youtubeService *youtube.Service

// InitYouTubeService initializes the YouTube Data API service
func InitYouTubeService() {
	yt, err := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		panic("Error intializing YouTube Data API Service: " + err.Error())
	}

	youtubeService = yt
}

// GetTitle returns the title of the video denoted by the given id
func GetTitle(id string) (string, error) {
	resp, err := youtubeService.Videos.List([]string{"snippet"}).Id(id).Do()
	if err != nil {
		return "", err
	}

	if len(resp.Items) == 0 {
		return "", nil
	}

	return resp.Items[0].Snippet.Title, nil
}
