package youtube_test

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/thohui/watchtogether/youtube"
)

func TestVideo(t *testing.T) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	service := youtube.New(apiKey)
	video, err := service.GetVideo("dQw4w9WgXcQ")
	if err != nil {
		t.Error(err)
	}
	if video.ID != "dQw4w9WgXcQ" {
		t.Error("invalid video id")
	}
	if video.Title != "Rick Astley - Never Gonna Give You Up (Official Music Video)" {
		t.Errorf("expected \"Rick Astley - Never Gonna Give You Up (Official Music Video)\", got \"%s\"", video.Title)
	}
	if video.Duration != 213 {
		t.Error("invalid video duration, expected 213, got", video.Duration)
	}
}
func TestDuration(t *testing.T) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	service := youtube.New(apiKey)
	duration := service.ParseDuration("PT1S")
	if duration != 1 {
		t.Error("invalid duration, expected 1, got", duration)
	}
	duration2 := service.ParseDuration("PT1M1S")
	if duration2 != 61 {
		t.Errorf("expected 61, got %d", duration2)
	}
	duration3 := service.ParseDuration("PT1H1M1S")
	if duration3 != 3661 {
		t.Errorf("expected 3661, got %d", duration3)
	}
}
