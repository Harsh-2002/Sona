package youtube

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DownloadAudio downloads audio from a YouTube URL using yt-dlp
func DownloadAudio(url string, outputDir string) (string, error) {
	fmt.Println("Downloading audio from YouTube...")

	// Check if yt-dlp is installed
	ytdlpPath, err := exec.LookPath("yt-dlp")
	if err != nil {
		// Try to install yt-dlp
		fmt.Println("yt-dlp not found, attempting to install...")
		if err := installYtDlp(); err != nil {
			return "", fmt.Errorf("failed to install yt-dlp: %v", err)
		}
		
		// Check again
		ytdlpPath, err = exec.LookPath("yt-dlp")
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
		fmt.Printf("Duration: %s\n", duration)
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

// installYtDlp attempts to install yt-dlp
func installYtDlp() error {
	// Try using pip first
	fmt.Println("Attempting to install yt-dlp using pip...")
	cmd := exec.Command("pip", "install", "--user", "yt-dlp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err == nil {
		return nil
	}
	
	// Try using pip3
	fmt.Println("Attempting to install yt-dlp using pip3...")
	cmd = exec.Command("pip3", "install", "--user", "yt-dlp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err == nil {
		return nil
	}
	
	// Try using curl to download directly
	fmt.Println("Attempting to download yt-dlp binary...")
	
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "yt-dlp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	tempFile.Close()
	
	// Download the binary
	cmd = exec.Command("curl", "-L", "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp", "-o", tempFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download yt-dlp: %v", err)
	}
	
	// Make it executable
	if err := os.Chmod(tempFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to make yt-dlp executable: %v", err)
	}
	
	// Try to move it to a location in PATH
	// First try user's home bin directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userBin := filepath.Join(homeDir, "bin")
		// Create bin directory if it doesn't exist
		if _, err := os.Stat(userBin); os.IsNotExist(err) {
			os.MkdirAll(userBin, 0755)
		}
		
		// Try to move to user's bin directory
		if err := os.Rename(tempFile.Name(), filepath.Join(userBin, "yt-dlp")); err == nil {
			// Add to PATH if not already there
			os.Setenv("PATH", os.Getenv("PATH")+":"+userBin)
			return nil
		}
	}
	
	// If user bin fails, try /usr/local/bin if we have permission
	if err := os.Rename(tempFile.Name(), "/usr/local/bin/yt-dlp"); err != nil {
		// If that fails too, keep it in the temp location and add to PATH
		os.Chmod(tempFile.Name(), 0755)
		newPath := filepath.Dir(tempFile.Name())
		os.Setenv("PATH", os.Getenv("PATH")+":"+newPath)
		fmt.Printf("Installed yt-dlp to: %s\n", tempFile.Name())
		return nil
	}
	
	fmt.Println("yt-dlp installed successfully")
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
	ytdlpPath, err := exec.LookPath("yt-dlp")
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
