package youtube

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	yt "github.com/kkdai/youtube/v2"
)

// DownloadAudio downloads audio from a YouTube URL using pure Go
// This makes the tool completely independent with no external dependencies
func DownloadAudio(url string, outputDir string) (string, error) {
	fmt.Println("📥 Downloading audio from YouTube using pure Go...")

	// Create a new YouTube client
	client := yt.Client{}

	// Get video info
	video, err := client.GetVideo(url)
	if err != nil {
		return "", fmt.Errorf("failed to get video info: %v", err)
	}

	fmt.Printf("🎬 Video: %s\n", video.Title)
	fmt.Printf("⏱️  Duration: %s\n", video.Duration)

	// Find the best audio format
	var audioFormat *yt.Format
	for _, format := range video.Formats {
		// Look for audio-only formats or formats with audio
		if format.AudioQuality != "" {
			if audioFormat == nil {
				audioFormat = &format
			} else if format.AudioQuality == "AUDIO_QUALITY_MEDIUM" || format.AudioQuality == "AUDIO_QUALITY_HIGH" {
				// Prefer higher quality audio
				audioFormat = &format
			}
		}
	}

	if audioFormat == nil {
		// Fallback: look for any format with audio
		for _, format := range video.Formats {
			if format.AudioQuality != "" {
				audioFormat = &format
				break
			}
		}
	}

	if audioFormat == nil {
		return "", fmt.Errorf("no audio format found for this video")
	}

	fmt.Printf("🎵 Audio format: %s\n", audioFormat.AudioQuality)

	// Generate output filename
	outputFile := filepath.Join(outputDir, "audio.mp4")
	
	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Download the audio stream
	fmt.Println("⬇️  Downloading audio stream...")
	stream, _, err := client.GetStream(video, audioFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get audio stream: %v", err)
	}
	defer stream.Close()

	// Copy the stream to the file
	_, err = io.Copy(file, stream)
	if err != nil {
		return "", fmt.Errorf("failed to save audio stream: %v", err)
	}

	fmt.Println("✅ Audio download completed successfully!")
	return outputFile, nil
}

// IsYouTubeURL checks if the given string is a YouTube URL
func IsYouTubeURL(url string) bool {
	return strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be")
}

// GetVideoInfo gets basic information about a YouTube video
func GetVideoInfo(url string) (*yt.Video, error) {
	client := yt.Client{}
	return client.GetVideo(url)
}
