package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// downloadYouTubeAudio downloads audio from a YouTube URL using yt-dlp
// yt-dlp is more efficient and actively maintained compared to youtube-dl
func downloadYouTubeAudio(url string, outputDir string) (string, error) {
	// Check if yt-dlp is installed
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp not found. Please install it first:\n" +
			"  - macOS: brew install yt-dlp\n" +
			"  - Ubuntu/Debian: sudo apt install yt-dlp\n" +
			"  - Or download from: https://github.com/yt-dlp/yt-dlp")
	}

	// Generate output filename
	outputFile := filepath.Join(outputDir, "audio.%(ext)s")

	// Build yt-dlp command with optimal settings for audio extraction
	args := []string{
		"--extract-audio",           // Extract audio only
		"--audio-format", "mp3",     // Convert to MP3 (widely supported)
		"--audio-quality", "0",      // Best quality
		"--no-playlist",             // Don't download playlists
		"--no-warnings",             // Reduce output noise
		"--output", outputFile,      // Output file pattern
		"--force-overwrites",        // Overwrite existing files
		url,                         // YouTube URL
	}

	// Execute yt-dlp command
	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("📥 Downloading audio from YouTube...")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %v", err)
	}

	// Find the downloaded file
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to read output directory: %v", err)
	}

	// Look for audio files
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasPrefix(name, "audio.") && (strings.HasSuffix(name, ".mp3") || 
				strings.HasSuffix(name, ".m4a") || strings.HasSuffix(name, ".webm")) {
				return filepath.Join(outputDir, name), nil
			}
		}
	}

	return "", fmt.Errorf("no audio file found after download")
}

// Alternative method using youtube-dl if yt-dlp is not available
func downloadYouTubeAudioFallback(url string, outputDir string) (string, error) {
	// Check if youtube-dl is installed
	if _, err := exec.LookPath("youtube-dl"); err != nil {
		return "", fmt.Errorf("neither yt-dlp nor youtube-dl found")
	}

	outputFile := filepath.Join(outputDir, "audio.%(ext)s")

	args := []string{
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--no-playlist",
		"--output", outputFile,
		"--force-overwrites",
		url,
	}

	cmd := exec.Command("youtube-dl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("📥 Downloading audio from YouTube (using youtube-dl)...")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("youtube-dl failed: %v", err)
	}

	// Find the downloaded file
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to read output directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasPrefix(name, "audio.") && (strings.HasSuffix(name, ".mp3") || 
				strings.HasSuffix(name, ".m4a") || strings.HasSuffix(name, ".webm")) {
				return filepath.Join(outputDir, name), nil
			}
		}
	}

	return "", fmt.Errorf("no audio file found after download")
}
