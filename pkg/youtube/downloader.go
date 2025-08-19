package youtube

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// DownloadAudio downloads audio from a YouTube URL using yt-dlp
func DownloadAudio(url string, outputDir string) (string, error) {
	fmt.Println("Downloading audio from YouTube...")

	// Check if yt-dlp is installed
	ytdlpPath, err := FindBinary("yt-dlp")
	if err != nil {
		// Try to install yt-dlp
		fmt.Println("yt-dlp not found, attempting to install...")
		if err := installYtDlp(); err != nil {
			return "", fmt.Errorf("failed to install yt-dlp: %v", err)
		}

		// Check again
		ytdlpPath, err = FindBinary("yt-dlp")
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

// FindBinary finds a binary in PATH or user's bin directory
func FindBinary(binaryName string) (string, error) {
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
	// Direct binary download is more reliable across platforms
	fmt.Println("Downloading yt-dlp binary directly...")
	return downloadYtDlpBinary()
}

// tryPipInstall attempts to install yt-dlp using pip (DEPRECATED)
func tryPipInstall() error {
	// This function is deprecated - direct binary download is preferred
	return fmt.Errorf("pip installation deprecated - using direct binary download")
}

// downloadYtDlpBinary downloads yt-dlp binary directly for the current platform
func downloadYtDlpBinary() error {
	// Determine platform and architecture
	platform := getPlatform()
	arch := getArchitecture()

	fmt.Printf("Detected platform: %s, architecture: %s\n", platform, arch)

	// Get the appropriate download URL for this platform
	downloadURL := getYtDlpDownloadURL(platform, arch)
	if downloadURL == "" {
		return fmt.Errorf("unsupported platform: %s/%s", platform, arch)
	}

	fmt.Printf("Download URL: %s\n", downloadURL)

	// Check if curl or wget is available
	var downloadCmd *exec.Cmd
	if _, err := exec.LookPath("curl"); err == nil {
		downloadCmd = exec.Command("curl", "-L", "-o", "yt-dlp", downloadURL)
	} else if _, err := exec.LookPath("wget"); err == nil {
		downloadCmd = exec.Command("wget", "-O", "yt-dlp", downloadURL)
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

	// Change to user's bin directory for download
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(userBin); err != nil {
		return fmt.Errorf("failed to change to bin directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Download the binary with verbose output
	fmt.Printf("Downloading to: %s\n", userBin)

	// Capture output for debugging
	var stderr bytes.Buffer
	downloadCmd.Stderr = &stderr

	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to download yt-dlp: %v\nStderr: %s", err, stderr.String())
	}

	// Verify the file was downloaded
	if _, err := os.Stat("yt-dlp"); err != nil {
		return fmt.Errorf("downloaded file not found: %v", err)
	}

	// Get file size
	if fileInfo, err := os.Stat("yt-dlp"); err == nil {
		fmt.Printf("Downloaded file size: %d bytes\n", fileInfo.Size())
		if fileInfo.Size() == 0 {
			return fmt.Errorf("downloaded file is empty")
		}
	}

	// Make it executable
	targetPath := filepath.Join(userBin, "yt-dlp")
	if err := os.Chmod(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to make yt-dlp executable: %v", err)
	}

	fmt.Printf("✅ yt-dlp installed successfully to: %s\n", targetPath)

	// Try to add to PATH for current session
	if err := addToPath(userBin); err != nil {
		fmt.Printf("⚠️  Warning: Could not update PATH. You may need to restart your terminal or run: export PATH=$PATH:%s\n", userBin)
	}

	return nil
}

// getPlatform returns the current platform
func getPlatform() string {
	switch runtime.GOOS {
	case "darwin":
		return "macos"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return runtime.GOOS
	}
}

// getArchitecture returns the current architecture
func getArchitecture() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	case "386":
		return "i386"
	default:
		return runtime.GOARCH
	}
}

// getYtDlpDownloadURL returns the appropriate download URL for the platform
func getYtDlpDownloadURL(platform, arch string) string {
	baseURL := "https://github.com/yt-dlp/yt-dlp/releases/latest/download"

	switch platform {
	case "macos":
		// macOS has universal binaries that work on both Intel and ARM64
		return baseURL + "/yt-dlp_macos"
	case "linux":
		if arch == "x86_64" {
			return baseURL + "/yt-dlp_linux"
		} else if arch == "aarch64" {
			return baseURL + "/yt-dlp_linux_aarch64"
		} else if arch == "armv7l" {
			return baseURL + "/yt-dlp_linux_armv7l"
		}
	case "windows":
		if arch == "x86_64" {
			return baseURL + "/yt-dlp.exe"
		} else if arch == "386" {
			return baseURL + "/yt-dlp_x86.exe"
		}
	}

	return ""
}

// addToPath attempts to add the bin directory to PATH for the current session
func addToPath(binDir string) error {
	// Get current PATH
	currentPath := os.Getenv("PATH")
	if currentPath == "" {
		currentPath = binDir
	} else {
		currentPath = binDir + ":" + currentPath
	}

	// Set PATH for current process
	return os.Setenv("PATH", currentPath)
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
	ytdlpPath, err := FindBinary("yt-dlp")
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
