package transcriber

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/root/sona-ai/pkg/assemblyai"
	"github.com/root/sona-ai/pkg/config"
	"github.com/root/sona-ai/pkg/youtube"
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

		// Get API key
		apiKey := config.GetAPIKey()
		fmt.Println("API key retrieved successfully")

		// Determine source type and process
		if youtube.IsYouTubeURL(source) {
			fmt.Println("Processing YouTube URL...")
			if err := processYouTubeVideo(source, apiKey); err != nil {
				fmt.Printf("Error: YouTube processing failed: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Processing local audio file...")
			if err := processLocalAudio(source, apiKey); err != nil {
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

func processYouTubeVideo(url string, apiKey string) error {
	// Get video info first
	title, _, err := youtube.GetVideoInfo(url)
	if err == nil {
		fmt.Printf("Processing: %s\n", title)
	}

	// Create temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "sona-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Download audio from YouTube
	fmt.Println("Downloading from YouTube...")

	audioPath, err := youtube.DownloadAudio(url, tempDir)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}

	// Transcribe the audio
	transcript, err := transcribeAudio(audioPath, apiKey)
	if err != nil {
		return fmt.Errorf("transcription failed: %v", err)
	}

	// Save transcript
	if err := saveTranscript(transcript, url, "youtube"); err != nil {
		return fmt.Errorf("failed to save transcript: %v", err)
	}

	return nil
}

func processLocalAudio(filePath string, apiKey string) error {
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
	transcript, err := transcribeAudio(convertedPath, apiKey)
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
	ffmpegPath, err := findFFmpegBinary()
	if err != nil {
		// Try to install ffmpeg
		fmt.Println("FFmpeg not found, attempting to install...")
		if err := installFFmpeg(); err != nil {
			return "", fmt.Errorf("FFmpeg is required for audio conversion. Please install it manually: %v", err)
		}

		// Check again
		ffmpegPath, err = findFFmpegBinary()
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

// findFFmpegBinary finds FFmpeg binary in PATH or user's bin directory
func findFFmpegBinary() (string, error) {
	// First check if it's in PATH
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}

	// Check user's bin directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userBinPath := filepath.Join(homeDir, "bin", "ffmpeg")
		if _, err := os.Stat(userBinPath); err == nil {
			return userBinPath, nil
		}
	}

	// Not found
	return "", fmt.Errorf("ffmpeg not found")
}

// installFFmpeg attempts to install FFmpeg
func installFFmpeg() error {
	// Try package managers first
	if err := tryPackageManagers(); err == nil {
		return nil
	}

	// If package managers fail, try direct download
	fmt.Println("Package managers not available or failed, trying direct download...")
	return downloadFFmpegBinary()
}

// tryPackageManagers attempts to install FFmpeg using various package managers
func tryPackageManagers() error {
	var cmd *exec.Cmd

	// Try apt-get (Debian/Ubuntu)
	if _, err := exec.LookPath("apt-get"); err == nil {
		fmt.Println("Attempting to install FFmpeg using apt-get...")
		cmd = exec.Command("sudo", "apt-get", "update")
		cmd.Run() // Ignore error

		cmd = exec.Command("sudo", "apt-get", "install", "-y", "ffmpeg")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try yum (CentOS/RHEL/Fedora)
	if _, err := exec.LookPath("yum"); err == nil {
		fmt.Println("Attempting to install FFmpeg using yum...")
		cmd = exec.Command("sudo", "yum", "install", "-y", "ffmpeg")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try dnf (newer Fedora)
	if _, err := exec.LookPath("dnf"); err == nil {
		fmt.Println("Attempting to install FFmpeg using dnf...")
		cmd = exec.Command("sudo", "dnf", "install", "-y", "ffmpeg")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try brew (macOS)
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Println("Attempting to install FFmpeg using brew...")
		cmd = exec.Command("brew", "install", "ffmpeg")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try choco (Windows)
	if _, err := exec.LookPath("choco"); err == nil {
		fmt.Println("Attempting to install FFmpeg using chocolatey...")
		cmd = exec.Command("choco", "install", "ffmpeg", "-y")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no suitable package manager found")
}

// downloadFFmpegBinary downloads FFmpeg binary for systems without package managers
func downloadFFmpegBinary() error {
	fmt.Println("Attempting to download FFmpeg binary...")

	// Get user's bin directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	userBin := filepath.Join(homeDir, "bin")
	if err := os.MkdirAll(userBin, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// For now, provide instructions rather than attempting binary download
	// (FFmpeg binary downloads are complex due to licensing and platform variations)
	fmt.Printf("\n")
	fmt.Printf("Unable to install FFmpeg automatically.\n")
	fmt.Printf("Please install FFmpeg manually:\n")
	fmt.Printf("\n")
	fmt.Printf("For macOS: brew install ffmpeg\n")
	fmt.Printf("For Ubuntu/Debian: sudo apt-get install ffmpeg\n")
	fmt.Printf("For CentOS/RHEL: sudo yum install ffmpeg\n")
	fmt.Printf("For Windows: Download from https://ffmpeg.org/download.html\n")
	fmt.Printf("\n")
	fmt.Printf("Or add FFmpeg to your PATH if already installed.\n")

	return fmt.Errorf("FFmpeg installation requires manual intervention")
}

func transcribeAudio(audioPath string, apiKey string) (string, error) {
	// Verify file exists
	_, err := os.Stat(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %v", err)
	}

	client := assemblyai.NewClient(apiKey)
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
			// Get video title from YouTube
			var err error
			title, _, err = youtube.GetVideoInfo(source)
			if err != nil {
				// Fallback to video ID if title can't be retrieved
				if strings.Contains(source, "v=") {
					parts := strings.Split(source, "v=")
					if len(parts) > 1 {
						videoID := strings.Split(parts[1], "&")[0]
						title = videoID
					}
				} else if strings.Contains(source, "youtu.be/") {
					parts := strings.Split(source, "youtu.be/")
					if len(parts) > 1 {
						videoID := strings.Split(parts[1], "?")[0]
						title = videoID
					}
				}
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
func ProcessYouTubeVideo(url string, apiKey string) error {
	return processYouTubeVideo(url, apiKey)
}

// ProcessLocalAudio processes a local audio file
func ProcessLocalAudio(filePath string, apiKey string) error {
	return processLocalAudio(filePath, apiKey)
}
