package youtube

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Harsh-2002/Sona/pkg/logger"
)

// DownloadAudio downloads audio from a YouTube URL using yt-dlp
func DownloadAudio(url string, outputDir string) (string, error) {
	logger.LogInfo("Downloading audio from YouTube URL: %s", url)

	// Check if yt-dlp is installed
	ytdlpPath, err := FindBinary("yt-dlp")
	if err != nil {
		logger.LogError("yt-dlp not found: %v", err)
		return "", fmt.Errorf("yt-dlp not found. Run 'sona install' to install dependencies")
	}

	logger.LogInfo("Using yt-dlp: %s", ytdlpPath)

	// Create output filename
	outputFilename := "youtube_audio.mp3"
	outputPath := filepath.Join(outputDir, outputFilename)

	// Get ffmpeg location for yt-dlp (consistent across Unix-like systems)
	ffmpegPath := ""

	// First try system PATH
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		ffmpegPath = path
	} else {
		// For Unix-like systems, try user's bin directory
		if runtime.GOOS != "windows" {
			homeDir, _ := os.UserHomeDir()
			userBinPath := filepath.Join(homeDir, "bin", "ffmpeg")
			if _, err := os.Stat(userBinPath); err == nil {
				ffmpegPath = userBinPath
			}
		}
	}

	// Build yt-dlp command with additional options for better compatibility
	args := []string{
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--output", outputPath,
		"--no-playlist",
		"--force-generic-extractor",
		"--extractor-args", "youtube:player_client=web",
	}

	// Add ffmpeg location if found
	if ffmpegPath != "" {
		args = append(args, "--ffmpeg-location", ffmpegPath)
		logger.LogInfo("Using ffmpeg at: %s", ffmpegPath)
	}

	args = append(args, url)

	logger.LogInfo("Running yt-dlp command: yt-dlp %v", args)

	// Execute yt-dlp
	cmd := exec.Command(ytdlpPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("yt-dlp command failed: %v, stderr: %s", err, stderr.String())

		// Try fallback options if first attempt fails
		logger.LogInfo("First attempt failed, trying fallback options")
		fallbackArgs := []string{
			"--extract-audio",
			"--audio-format", "mp3",
			"--audio-quality", "0",
			"--output", outputPath,
			"--no-playlist",
			"--extractor-args", "youtube:player_client=web",
		}

		// Add ffmpeg location to fallback as well
		if ffmpegPath != "" {
			fallbackArgs = append(fallbackArgs, "--ffmpeg-location", ffmpegPath)
		}

		fallbackArgs = append(fallbackArgs, url)

		cmd = exec.Command(ytdlpPath, fallbackArgs...)
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			logger.LogError("yt-dlp fallback also failed: %v, stderr: %s", err, stderr.String())
			return "", fmt.Errorf("failed to download audio: %v", err)
		}

		logger.LogInfo("yt-dlp fallback succeeded")
	}

	logger.LogInfo("Audio download completed successfully: %s", outputPath)
	return outputPath, nil
}

// FindBinary finds a binary in PATH or user's bin directory
func FindBinary(binaryName string) (string, error) {
	// First check if it's in PATH
	if path, err := exec.LookPath(binaryName); err == nil {
		return path, nil
	}

	// For Unix-like systems, check user's bin directory
	if runtime.GOOS != "windows" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			userBinPath := filepath.Join(homeDir, "bin", binaryName)
			if _, err := os.Stat(userBinPath); err == nil {
				return userBinPath, nil
			}
		}
	}

	// Not found
	return "", fmt.Errorf("%s not found", binaryName)
}

// InstallYtDlp attempts to install yt-dlp
func InstallYtDlp() error {
	// Direct binary download is more reliable across platforms
	logger.LogInfo("Installing yt-dlp binary directly")
	return downloadYtDlpBinary()
}

// downloadYtDlpBinary downloads yt-dlp binary directly for the current platform
func downloadYtDlpBinary() error {
	platform, arch := getPlatform(), getArchitecture()
	logger.LogInfo("Detected platform: %s, architecture: %s", platform, arch)

	downloadURL := getYtDlpDownloadURL(platform, arch)
	if downloadURL == "" {
		return fmt.Errorf("unsupported platform: %s-%s", platform, arch)
	}

	logger.LogInfo("Download URL: %s", downloadURL)

	// Create bin directory if it doesn't exist (consistent path across Unix-like systems)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	binDir := filepath.Join(homeDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Download the binary
	outputPath := filepath.Join(binDir, "yt-dlp")
	logger.LogInfo("Downloading yt-dlp binary to: %s", binDir)

	cmd := exec.Command("curl", "-L", "-o", outputPath, downloadURL)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.LogError("Failed to download yt-dlp: %v, output: %s", err, string(output))
		return fmt.Errorf("download failed: %v", err)
	}

	// Make it executable
	if err := os.Chmod(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to make yt-dlp executable: %v", err)
	}

	// Verify the download
	if info, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("failed to verify download: %v", err)
	} else {
		logger.LogInfo("Downloaded file size: %d bytes", info.Size())
	}

	logger.LogInfo("yt-dlp installed successfully to: %s", outputPath)
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
