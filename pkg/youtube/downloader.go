package youtube

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// DownloadAudio downloads audio from a YouTube URL using yt-dlp
func DownloadAudio(url string, outputDir string) (string, error) {
	fmt.Println("Downloading audio from YouTube...")

	// Check if yt-dlp is installed
	ytdlpPath, err := findBinary("yt-dlp")
	if err != nil {
		// Try to install yt-dlp
		fmt.Println("yt-dlp not found, attempting to install...")
		if err := installYtDlp(); err != nil {
			return "", fmt.Errorf("failed to install yt-dlp: %v", err)
		}

		// Check again
		ytdlpPath, err = findBinary("yt-dlp")
		if err != nil {
			return "", fmt.Errorf("yt-dlp not found after installation attempt: %v", err)
		}
	}

	fmt.Printf("Using yt-dlp: %s\n", ytdlpPath)

	// Create output filename
	outputFile := filepath.Join(outputDir, "audio.mp3")

	// Get video info first
	title, duration, err := getVideoInfo(url)
	if err != nil {
		fmt.Printf("Warning: Could not get video info: %v\n", err)
		title = "Unknown"
		duration = "Unknown"
	} else {
		fmt.Printf("Video: %s\n", title)
		fmt.Printf("Duration: %s\n", formatDuration(duration))
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Download the audio using yt-dlp
	fmt.Println("Downloading audio stream...")

	// Prepare command
	cmd := exec.CommandContext(ctx, ytdlpPath,
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--output", outputFile,
		"--no-playlist",
		"--progress",
		"--quiet",
		url,
	)

	// Redirect output to null to hide technical details
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to download audio: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputFile); err != nil {
		return "", fmt.Errorf("failed to find downloaded file: %v", err)
	}

	// Get file size
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	fmt.Printf("Audio download completed successfully (%.2f MB)\n", float64(fileInfo.Size())/1024/1024)

	return outputFile, nil
}

// formatDuration converts duration string to human-readable format
// Examples: "1342" -> "22m 22s", "90" -> "1m 30s", "45" -> "45s"
func formatDuration(duration string) string {
	// Try to parse as integer seconds
	if seconds, err := strconv.Atoi(duration); err == nil {
		minutes := seconds / 60
		remainingSeconds := seconds % 60
		
		if minutes > 0 {
			if remainingSeconds > 0 {
				return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
			}
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%ds", remainingSeconds)
	}
	
	// If parsing fails, return as-is (might already be formatted)
	return duration
}

// findBinary finds a binary in PATH or user's bin directory
func findBinary(binaryName string) (string, error) {
	// First check if it's in PATH
	if path, err := exec.LookPath(binaryName); err == nil {
		return path, nil
	}

	// Check user's bin directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userBinPath := filepath.Join(homeDir, "bin", binaryName)
		if _, err := os.Stat(userBinPath); err == nil {
			return userBinPath, nil
		}
	}

	// Not found
	return "", fmt.Errorf("%s not found", binaryName)
}

// installYtDlp attempts to install yt-dlp
func installYtDlp() error {
	// Try pip installations first
	if err := tryPipInstall(); err == nil {
		return nil
	}

	// If pip fails, try direct binary download
	fmt.Println("pip installation failed, trying direct download...")
	return downloadYtDlpBinary()
}

// tryPipInstall attempts to install yt-dlp using pip
func tryPipInstall() error {
	// Try using pip first
	if _, err := exec.LookPath("pip"); err == nil {
		fmt.Println("Attempting to install yt-dlp using pip...")
		cmd := exec.Command("pip", "install", "--user", "yt-dlp")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try using pip3
	if _, err := exec.LookPath("pip3"); err == nil {
		fmt.Println("Attempting to install yt-dlp using pip3...")
		cmd := exec.Command("pip3", "install", "--user", "yt-dlp")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("pip installation failed")
}

// downloadYtDlpBinary downloads yt-dlp binary directly
func downloadYtDlpBinary() error {
	// Check if curl or wget is available
	var downloadCmd *exec.Cmd
	if _, err := exec.LookPath("curl"); err == nil {
		downloadCmd = exec.Command("curl", "-L",
			"https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp")
	} else if _, err := exec.LookPath("wget"); err == nil {
		downloadCmd = exec.Command("wget", "-O", "-",
			"https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp")
	} else {
		return fmt.Errorf("neither curl nor wget found - cannot download yt-dlp")
	}

	fmt.Println("Downloading yt-dlp binary...")

	// Get user's bin directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	userBin := filepath.Join(homeDir, "bin")
	if err := os.MkdirAll(userBin, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Download directly to target location
	targetPath := filepath.Join(userBin, "yt-dlp")

	if downloadCmd.Args[0] == "curl" {
		downloadCmd.Args = append(downloadCmd.Args, "-o", targetPath)
	} else {
		// For wget, we already set -O -
		outputFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create target file: %v", err)
		}
		defer outputFile.Close()
		downloadCmd.Stdout = outputFile
	}

	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to download yt-dlp: %v", err)
	}

	// Make it executable
	if err := os.Chmod(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to make yt-dlp executable: %v", err)
	}

	fmt.Printf("✅ yt-dlp installed successfully to: %s\n", targetPath)
	return nil
}

// IsYouTubeURL checks if the given string is a YouTube URL
func IsYouTubeURL(url string) bool {
	return strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be")
}

// getVideoInfo gets basic information about a YouTube video
func getVideoInfo(url string) (string, string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if yt-dlp is installed
	ytdlpPath, err := findBinary("yt-dlp")
	if err != nil {
		return "", "", fmt.Errorf("yt-dlp not found: %v", err)
	}

	// Get video info using yt-dlp
	cmd := exec.CommandContext(ctx, ytdlpPath,
		"--print", "title",
		"--print", "duration",
		"--no-playlist",
		"--no-download",
		url,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get video info: %v", err)
	}

	// Parse output (title and duration are on separate lines)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return "", "", fmt.Errorf("unexpected output format from yt-dlp")
	}

	title := lines[0]
	duration := lines[1]

	return title, duration, nil
}

// GetVideoInfo gets basic information about a YouTube video (public API)
func GetVideoInfo(url string) (string, string, error) {
	return getVideoInfo(url)
}
