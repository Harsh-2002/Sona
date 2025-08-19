package transcriber

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Harsh-2002/Sona/pkg/assemblyai"
	"github.com/Harsh-2002/Sona/pkg/config"
	"github.com/Harsh-2002/Sona/pkg/logger"
	"github.com/Harsh-2002/Sona/pkg/youtube"
	"github.com/spf13/cobra"
)

var (
	outputPath  string
	speechModel string
)

var TranscribeCmd = &cobra.Command{
	Use:   "transcribe [source]",
	Short: "Transcribe audio from YouTube video or local file",
	Long: `Transcribe audio to text using AssemblyAI.
	
Sources:
- YouTube URL: sona transcribe "https://youtube.com/watch?v=..."
- Local file: sona transcribe "./audio.mp3"

Examples:
  sona transcribe "https://youtube.com/watch?v=dQw4w9WgXcQ"
  sona transcribe "./audio.mp3"
  sona transcribe "https://youtube.com/watch?v=..." --output ./transcript.txt
  sona transcribe "./audio.mp3" --model slam-1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		fmt.Printf("Source: %s\n", source)

		// Check and install dependencies
		if err := checkAndInstallDependencies(); err != nil {
			fmt.Printf("Error: Dependency check failed: %v\n", err)
			os.Exit(1)
		}

		// Determine source type and process
		if youtube.IsYouTubeURL(source) {
			fmt.Println("Processing YouTube URL...")
			if err := processYouTubeVideo(source, outputPath, speechModel); err != nil {
				fmt.Printf("Error: YouTube processing failed: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Processing local audio file...")
			if err := processLocalAudio(source, outputPath, speechModel); err != nil {
				fmt.Printf("Error: Local audio processing failed: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("Transcription completed successfully")
	},
}

func init() {
	TranscribeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: auto-generated)")
	TranscribeCmd.Flags().StringVarP(&speechModel, "model", "m", "slam-1", "Speech model to use (slam-1, best, nano)")
}

// checkAndInstallDependencies ensures both yt-dlp and ffmpeg are available
func checkAndInstallDependencies() error {
	fmt.Println("🔍 Checking dependencies...")
	logger.LogInfo("Checking and installing dependencies")

	// Check yt-dlp
	ytdlpPath, err := youtube.FindBinary("yt-dlp")
	if err != nil {
		fmt.Println("📥 Installing yt-dlp...")
		logger.LogInfo("yt-dlp not found, installing")
		if err := youtube.InstallYtDlp(); err != nil {
			logger.LogError("Failed to install yt-dlp: %v", err)
			return fmt.Errorf("failed to install yt-dlp: %v", err)
		}
		fmt.Println("✅ yt-dlp installed successfully")
	} else {
		logger.LogInfo("yt-dlp found at: %s", ytdlpPath)
	}

	// Check ffmpeg
	ffmpegPath, err := FindBinary("ffmpeg")
	if err != nil {
		fmt.Println("📥 Installing FFmpeg...")
		logger.LogInfo("FFmpeg not found, installing")
		if err := installFFmpeg(); err != nil {
			logger.LogError("Failed to install FFmpeg: %v", err)
			return fmt.Errorf("failed to install FFmpeg: %v", err)
		}
		fmt.Println("✅ FFmpeg installed successfully")
	} else {
		logger.LogInfo("FFmpeg found at: %s", ffmpegPath)

		// On macOS, also check for ffprobe
		if runtime.GOOS == "darwin" {
			if _, err := FindBinary("ffprobe"); err != nil {
				fmt.Println("📥 Installing ffprobe for macOS...")
				logger.LogInfo("ffprobe not found on macOS, installing")
				if err := installFFmpeg(); err != nil {
					logger.LogError("Failed to install ffprobe: %v", err)
					return fmt.Errorf("failed to install ffprobe: %v", err)
				}
				fmt.Println("✅ ffprobe installed successfully")
			} else {
				logger.LogInfo("ffprobe found")
			}
		}
	}

	fmt.Println("🎯 All dependencies are ready!")
	return nil
}

func processYouTubeVideo(url string, outputPath string, speechModel string) error {
	fmt.Println("Processing YouTube URL...")
	logger.LogInfo("Processing YouTube video: %s", url)

	// Download audio from YouTube
	audioFile, err := youtube.DownloadAudio(url, filepath.Dir(outputPath))
	if err != nil {
		logger.LogError("Failed to download YouTube audio: %v", err)
		return fmt.Errorf("failed to download YouTube audio: %v", err)
	}

	logger.LogInfo("Audio downloaded successfully: %s", audioFile)

	// Transcribe the audio
	transcript, err := transcribeAudio(audioFile, speechModel)
	if err != nil {
		logger.LogError("Failed to transcribe YouTube audio: %v", err)
		return fmt.Errorf("failed to transcribe audio: %v", err)
	}

	// Save transcript
	if err := saveTranscript(transcript, url, "youtube"); err != nil {
		logger.LogError("Failed to save transcript: %v", err)
		return fmt.Errorf("failed to save transcript: %v", err)
	}

	// Clean up audio file
	os.Remove(audioFile)
	logger.LogInfo("YouTube video processing completed successfully")

	return nil
}

func processLocalAudio(filePath string, outputPath string, speechModel string) error {
	// Check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("audio file not found: %s", filePath)
	}

	// Show file info
	fmt.Printf("Processing: %s\n", filepath.Base(filePath))

	// Create temporary directory for conversion
	tempDir, err := os.MkdirTemp("", "sona-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Convert audio to MP3 format for better compatibility
	convertedPath, err := convertAudioToMP3(filePath, tempDir)
	if err != nil {
		return fmt.Errorf("audio conversion failed: %v", err)
	}

	// Transcribe the converted audio
	transcript, err := transcribeAudio(convertedPath, speechModel)
	if err != nil {
		return fmt.Errorf("transcription failed: %v", err)
	}

	// Save transcript
	if err := saveTranscript(transcript, filePath, "local"); err != nil {
		return fmt.Errorf("failed to save transcript: %v", err)
	}

	return nil
}

// convertAudioToMP3 converts audio file to MP3 format for better compatibility
func convertAudioToMP3(inputPath string, outputDir string) (string, error) {
	// Check if ffmpeg is installed
	ffmpegPath, err := FindBinary("ffmpeg")
	if err != nil {
		// Try to install ffmpeg
		fmt.Println("FFmpeg not found, attempting to install...")
		if err := installFFmpeg(); err != nil {
			return "", fmt.Errorf("FFmpeg is required for audio conversion. Please install it manually: %v", err)
		}

		// Check again
		ffmpegPath, err = FindBinary("ffmpeg")
		if err != nil {
			return "", fmt.Errorf("FFmpeg not found after installation attempt: %v", err)
		}
	}

	// Create output path
	outputPath := filepath.Join(outputDir, "converted.mp3")

	fmt.Println("Converting audio to MP3 format...")

	// Run ffmpeg to convert the file
	cmd := exec.Command(ffmpegPath,
		"-i", inputPath,
		"-vn",          // No video
		"-ar", "44100", // Sample rate
		"-ac", "2", // Stereo
		"-b:a", "192k", // Bitrate
		"-f", "mp3", // Format
		"-y", // Overwrite output
		outputPath)

	// Hide ffmpeg output
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to convert audio: %v", err)
	}

	// Verify the converted file exists
	if _, err := os.Stat(outputPath); err != nil {
		return "", fmt.Errorf("converted file not found: %v", err)
	}

	fmt.Println("Audio conversion completed")
	return outputPath, nil
}

// FindBinary finds FFmpeg binary in PATH or user's bin directory
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

// installFFmpeg attempts to install FFmpeg
func installFFmpeg() error {
	// Direct binary download is more reliable across platforms
	fmt.Println("Downloading FFmpeg binary directly...")
	return downloadFFmpegBinary()
}

// downloadFFmpegBinary downloads FFmpeg binary directly for the current platform
func downloadFFmpegBinary() error {
	fmt.Println("Attempting to download FFmpeg binary...")

	platform := getPlatform()
	arch := getArchitecture()

	logger.LogInfo("Detected platform: %s, architecture: %s", platform, arch)

	if platform == "macos" {
		// For macOS, download both ffmpeg and ffprobe from evermeet.cx
		return downloadMacOSFFmpeg()
	}

	// For other platforms, use BtbN builds
	downloadURL, filename := getFFmpegDownloadURL(platform, arch)
	if downloadURL == "" {
		return fmt.Errorf("unsupported platform: %s/%s", platform, arch)
	}

	logger.LogInfo("Downloading FFmpeg from: %s", downloadURL)

	// Create bin directory if it doesn't exist
	binDir := filepath.Join(os.Getenv("HOME"), "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Change to the bin directory for extraction
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(binDir); err != nil {
		return fmt.Errorf("failed to change to bin directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Download the archive
	cmd := exec.Command("curl", "-L", "-o", filename, downloadURL)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download FFmpeg: %v", err)
	}

	// Extract the archive
	if err := extractFFmpegArchive(filename); err != nil {
		return fmt.Errorf("failed to extract FFmpeg archive: %v", err)
	}

	logger.LogInfo("FFmpeg installed successfully")
	return nil
}

// downloadMacOSFFmpeg downloads ffmpeg and ffprobe for macOS from evermeet.cx
func downloadMacOSFFmpeg() error {
	logger.LogInfo("Downloading FFmpeg and ffprobe for macOS from evermeet.cx")

	binDir := filepath.Join(os.Getenv("HOME"), "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Download ffmpeg
	ffmpegURL := "https://evermeet.cx/ffmpeg/get/zip"
	ffmpegPath := filepath.Join(binDir, "ffmpeg.zip")
	logger.LogInfo("Downloading ffmpeg from: %s", ffmpegURL)

	cmd := exec.Command("curl", "-L", "-o", ffmpegPath, ffmpegURL)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.LogError("Failed to download ffmpeg: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to download ffmpeg: %v", err)
	}

	// Extract ffmpeg
	cmd = exec.Command("unzip", "-q", "-o", ffmpegPath, "-d", binDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.LogError("Failed to extract ffmpeg: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to extract ffmpeg: %v", err)
	}

	// Download ffprobe
	ffprobeURL := "https://evermeet.cx/ffmpeg/get/ffprobe/zip"
	ffprobePath := filepath.Join(binDir, "ffprobe.zip")
	logger.LogInfo("Downloading ffprobe from: %s", ffprobeURL)

	cmd = exec.Command("curl", "-L", "-o", ffprobePath, ffprobeURL)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.LogError("Failed to download ffprobe: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to download ffprobe: %v", err)
	}

	// Extract ffprobe
	cmd = exec.Command("unzip", "-q", "-o", ffprobePath, "-d", binDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.LogError("Failed to extract ffprobe: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to extract ffprobe: %v", err)
	}

	// Make binaries executable
	ffmpegBin := filepath.Join(binDir, "ffmpeg")
	ffprobeBin := filepath.Join(binDir, "ffprobe")

	if err := os.Chmod(ffmpegBin, 0755); err != nil {
		return fmt.Errorf("failed to make ffmpeg executable: %v", err)
	}

	if err := os.Chmod(ffprobeBin, 0755); err != nil {
		return fmt.Errorf("failed to make ffprobe executable: %v", err)
	}

	// Clean up zip files
	os.Remove(ffmpegPath)
	os.Remove(ffprobePath)

	logger.LogInfo("FFmpeg and ffprobe installed successfully to: %s", binDir)
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

// getFFmpegDownloadURL returns the appropriate download URL and filename for the platform
func getFFmpegDownloadURL(platform, arch string) (string, string) {
	switch platform {
	case "macos":
		if arch == "x86_64" {
			// Use evermeet.cx for macOS Intel (more reliable)
			return "https://evermeet.cx/ffmpeg/ffmpeg-120751-g1d06e8ddcd.zip", "ffmpeg-macos-intel.zip"
		} else if arch == "aarch64" {
			// Use evermeet.cx for macOS ARM64 (more reliable)
			return "https://evermeet.cx/ffmpeg/ffmpeg-120751-g1d06e8ddcd.zip", "ffmpeg-macos-arm64.zip"
		}
	case "linux":
		if arch == "x86_64" {
			// Use static builds from BtbN's repository for Linux
			return "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz", "ffmpeg-linux64.tar.xz"
		} else if arch == "aarch64" {
			return "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linuxarm64-gpl.tar.xz", "ffmpeg-linuxarm64.tar.xz"
		}
	case "windows":
		if arch == "x86_64" {
			// Use static builds from BtbN's repository for Windows
			return "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip", "ffmpeg-win64.zip"
		} else if arch == "aarch64" {
			return "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-winarm64-gpl.zip", "ffmpeg-winarm64.zip"
		}
	}

	return "", ""
}

// extractFFmpegArchive extracts the downloaded FFmpeg archive
func extractFFmpegArchive(filename string) error {
	fmt.Printf("Extracting %s...\n", filename)

	var cmd *exec.Cmd

	if strings.HasSuffix(filename, ".tar.gz") {
		cmd = exec.Command("tar", "-xzf", filename)
	} else if strings.HasSuffix(filename, ".tar.xz") {
		cmd = exec.Command("tar", "-xf", filename)
	} else if strings.HasSuffix(filename, ".zip") {
		cmd = exec.Command("unzip", "-q", filename)
	} else {
		return fmt.Errorf("unsupported archive format: %s", filename)
	}

	// Capture stderr for better error reporting
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract archive: %v\nStderr: %s", err, stderr.String())
	}

	// Find the ffmpeg binary in the extracted directory
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Look for the ffmpeg binary
	var ffmpegFound bool
	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), "ffmpeg") {
			// Check if there's a bin subdirectory
			binPath := filepath.Join(entry.Name(), "bin", "ffmpeg")
			if _, err := os.Stat(binPath); err == nil {
				// Move the binary to the user's bin directory
				finalPath := filepath.Join(".", "ffmpeg")
				if err := os.Rename(binPath, finalPath); err != nil {
					return fmt.Errorf("failed to move FFmpeg binary: %v", err)
				}
				ffmpegFound = true
				break
			}
		}
	}

	// For macOS ZIP files, the binary might be directly in the archive
	if !ffmpegFound {
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() == "ffmpeg" {
				// Binary is already in the right place
				ffmpegFound = true
				break
			}
		}
	}

	if !ffmpegFound {
		// List what we found for debugging
		fmt.Println("Debug: Found entries after extraction:")
		for _, entry := range entries {
			fmt.Printf("  - %s (dir: %t)\n", entry.Name(), entry.IsDir())
		}
		return fmt.Errorf("could not find FFmpeg binary in extracted archive")
	}

	// Clean up extracted files and archive
	for _, entry := range entries {
		if entry.IsDir() {
			os.RemoveAll(entry.Name())
		}
	}
	os.Remove(filename)

	return nil
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

func transcribeAudio(audioPath string, speechModel string) (string, error) {
	// Verify file exists
	_, err := os.Stat(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %v", err)
	}

	client := assemblyai.NewClient(config.GetAPIKey())
	return client.TranscribeAudio(audioPath, speechModel)
}

func saveTranscript(transcript string, source string, sourceType string) error {
	// Determine output path
	var finalOutputPath string
	if outputPath != "" {
		finalOutputPath = outputPath
	} else {
		// Generate default path
		defaultPath := config.GetOutputPath()
		if err := os.MkdirAll(defaultPath, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}

		// Generate filename based on source
		var filename string
		var title string

		if sourceType == "youtube" {
			// Extract video ID from YouTube URL for filename
			if strings.Contains(source, "v=") {
				parts := strings.Split(source, "v=")
				if len(parts) > 1 {
					videoID := strings.Split(parts[1], "&")[0]
					title = "youtube-" + videoID
				}
			} else if strings.Contains(source, "youtu.be/") {
				parts := strings.Split(source, "youtu.be/")
				if len(parts) > 1 {
					videoID := strings.Split(parts[1], "?")[0]
					title = "youtube-" + videoID
				}
			}
			if title == "" {
				title = "youtube-video"
			}
		} else {
			// For local files, use the filename without extension
			baseName := filepath.Base(source)
			ext := filepath.Ext(baseName)
			if len(ext) > 0 && len(baseName) > len(ext) {
				title = baseName[:len(baseName)-len(ext)]
			} else {
				title = baseName
			}
		}

		// Sanitize title for use as filename
		title = sanitizeFilename(title)

		// If title is empty or couldn't be determined, use a default
		if title == "" {
			title = "transcript"
		}

		// Add simple timestamp for uniqueness (just date)
		timestamp := time.Now().Format("20060102")
		filename = fmt.Sprintf("%s-%s.txt", title, timestamp)

		finalOutputPath = filepath.Join(defaultPath, filename)
	}

	// Write transcript to file
	if err := os.WriteFile(finalOutputPath, []byte(transcript), 0644); err != nil {
		return fmt.Errorf("failed to write transcript file: %v", err)
	}

	fmt.Printf("Saved to: %s (%d chars)\n", finalOutputPath, len(transcript))

	return nil
}

// sanitizeFilename removes invalid characters from a filename and makes it cleaner
func sanitizeFilename(name string) string {
	// Replace invalid characters with hyphens
	reg := regexp.MustCompile(`[\\/:*?"<>|]`)
	name = reg.ReplaceAllString(name, "-")

	// Replace spaces and underscores with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Replace multiple hyphens with a single hyphen
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	// Remove leading/trailing spaces and hyphens
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "-")

	// Convert to lowercase for consistency
	name = strings.ToLower(name)

	// Limit length to avoid too long filenames
	const maxLength = 40
	if len(name) > maxLength {
		name = name[:maxLength]
	}

	// Ensure name is not empty
	if name == "" {
		name = "transcript"
	}

	return name
}

// SetOutputPath sets the output path for the transcript
func SetOutputPath(path string) {
	outputPath = path
}

// SetSpeechModel sets the speech model to use
func SetSpeechModel(model string) {
	speechModel = model
}

// ProcessYouTubeVideo processes a YouTube video URL
func ProcessYouTubeVideo(url string, outputPath string, speechModel string) error {
	return processYouTubeVideo(url, outputPath, speechModel)
}

// ProcessLocalAudio processes a local audio file
func ProcessLocalAudio(filePath string, outputPath string, speechModel string) error {
	return processLocalAudio(filePath, outputPath, speechModel)
}
