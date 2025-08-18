package transcriber

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/root/sona-ai/pkg/assemblyai"
	"github.com/root/sona-ai/pkg/config"
	"github.com/root/sona-ai/pkg/youtube"
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
		
		// Get API key
		apiKey := config.GetAPIKey()
		
		// Determine source type and process
		if youtube.IsYouTubeURL(source) {
			fmt.Println("🎥 Detected YouTube URL, downloading audio...")
			if err := processYouTubeVideo(source, apiKey); err != nil {
				fmt.Printf("❌ Error processing YouTube video: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("🎵 Processing local audio file...")
			if err := processLocalAudio(source, apiKey); err != nil {
				fmt.Printf("❌ Error processing local audio: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	TranscribeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: auto-generated)")
	TranscribeCmd.Flags().StringVarP(&speechModel, "model", "m", "slam-1", "Speech model to use (slam-1, best, nano)")
}

func processYouTubeVideo(url string, apiKey string) error {
	// Create temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "sona-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Download audio from YouTube using pure Go (no external dependencies)
	audioPath, err := youtube.DownloadAudio(url, tempDir)
	if err != nil {
		return fmt.Errorf("failed to download YouTube audio: %v", err)
	}

	fmt.Printf("✅ Downloaded audio to: %s\n", audioPath)

	// Transcribe the audio
	transcript, err := transcribeAudio(audioPath, apiKey)
	if err != nil {
		return fmt.Errorf("failed to transcribe audio: %v", err)
	}

	// Save transcript
	if err := saveTranscript(transcript, url, "youtube"); err != nil {
		return fmt.Errorf("failed to save transcript: %v", err)
	}

	return nil
}

func processLocalAudio(filePath string, apiKey string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("audio file not found: %s", filePath)
	}

	// Transcribe the audio
	transcript, err := transcribeAudio(filePath, apiKey)
	if err != nil {
		return fmt.Errorf("failed to transcribe audio: %v", err)
	}

	// Save transcript
	if err := saveTranscript(transcript, filePath, "local"); err != nil {
		return fmt.Errorf("failed to save transcript: %v", err)
	}

	return nil
}

func transcribeAudio(audioPath string, apiKey string) (string, error) {
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
		if sourceType == "youtube" {
			// Extract video ID from YouTube URL
			if strings.Contains(source, "v=") {
				parts := strings.Split(source, "v=")
				if len(parts) > 1 {
					videoID := strings.Split(parts[1], "&")[0]
					filename = fmt.Sprintf("youtube_%s.txt", videoID)
				}
			}
		} else {
			// Use original filename with .txt extension
			baseName := filepath.Base(source)
			ext := filepath.Ext(baseName)
			filename = baseName[:len(baseName)-len(ext)] + "_transcript.txt"
		}

		if filename == "" {
			filename = "transcript.txt"
		}

		finalOutputPath = filepath.Join(defaultPath, filename)
	}

	// Write transcript to file
	if err := os.WriteFile(finalOutputPath, []byte(transcript), 0644); err != nil {
		return fmt.Errorf("failed to write transcript file: %v", err)
	}

	fmt.Printf("✅ Transcript saved to: %s\n", finalOutputPath)
	fmt.Printf("📝 Transcript length: %d characters\n", len(transcript))
	
	return nil
}
