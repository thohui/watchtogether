package youtube

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeService struct {
	service *youtube.Service
}

type Video struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

func New(apiKey string) *YoutubeService {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}
	youtubeService := &YoutubeService{service: service}
	return youtubeService
}

func (service *YoutubeService) GetVideo(id string) (Video, error) {
	if len(id) != 11 {
		return Video{}, errors.New("invalid video id")
	}
	var args = []string{"id", "snippet", "contentDetails"}
	call := service.service.Videos.List(args).Id(id)
	response, err := call.Do()
	if err != nil {
		return Video{}, err
	}
	if len(response.Items) == 0 {
		return Video{}, errors.New("video not found")
	}
	video := response.Items[0]
	duration := service.ParseDuration(video.ContentDetails.Duration)
	return Video{
		ID:       video.Id,
		Title:    video.Snippet.Title,
		Duration: duration,
	}, nil
}

func (service *YoutubeService) ParseDuration(duration string) int {
	result := 0
	duration = strings.Replace(duration, "PT", "", 1)
	if strings.Contains(duration, "H") {
		split := strings.Split(duration, "H")
		hours, err := strconv.Atoi(split[0])
		if err != nil {
			return 0
		}
		length := len(split[0]) + 1
		duration = duration[length:]
		result += hours * 60 * 60
	}
	if strings.Contains(duration, "M") {
		split := strings.Split(duration, "M")
		minutes, err := strconv.Atoi(split[0])
		if err != nil {
			return 0
		}
		length := len(split[0]) + 1
		duration = duration[length:]
		result += minutes * 60
	}
	if strings.Contains(duration, "S") {
		split := strings.Split(duration, "S")
		seconds, err := strconv.Atoi(split[0])
		if err == nil {
			result += seconds
		}

	}
	return result
}
