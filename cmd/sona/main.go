package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Harsh-2002/Sona/pkg/config"
	"github.com/Harsh-2002/Sona/pkg/interactive"
	"github.com/Harsh-2002/Sona/pkg/logger"
	"github.com/Harsh-2002/Sona/pkg/transcriber"
	"github.com/Harsh-2002/Sona/pkg/youtube"
	"github.com/spf13/cobra"
)

// Version will be set by the build process
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "sona",
	Short: "Audio Transcription Tool",
	Long: `Sona - Audio Transcription Tool

A CLI tool that converts audio files and YouTube videos to text transcripts using AssemblyAI.

Features:
- Transcribe local audio files
- Download and transcribe YouTube videos
- Save transcripts to custom or default paths
- Interactive mode for guided experience`,
	Run: func(cmd *cobra.Command, args []string) {
		interactive.InteractiveCmd.Run(cmd, args)
	},
}

func init() {
	// Initialize configuration
	config.InitConfig()

	// Add commands
	rootCmd.AddCommand(transcriber.TranscribeCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(interactive.InteractiveCmd)
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check system status and dependencies",
	Long:  "Check the status of yt-dlp and FFmpeg dependencies and system configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sona System Status")
		fmt.Println("==================")

		// Check yt-dlp
		fmt.Println("\n1. YouTube Download (yt-dlp):")
		if ytdlpPath, err := youtube.FindBinary("yt-dlp"); err == nil {
			fmt.Printf("   Available at: %s\n", ytdlpPath)
		} else {
			fmt.Println("   Not found (will auto-install when needed)")
		}

		// Check FFmpeg
		fmt.Println("\n2. Audio Processing (FFmpeg):")
		if ffmpegPath, err := transcriber.FindBinary("ffmpeg"); err == nil {
			fmt.Printf("   Available at: %s\n", ffmpegPath)
		} else {
			fmt.Println("   Not found (will auto-install when needed)")
		}

		// Check API key
		fmt.Println("\n3. AssemblyAI API Key:")
		apiKey := config.GetAPIKeyNoExit()
		if apiKey != "" {
			fmt.Println("   Configured")
		} else {
			fmt.Println("   Not configured")
			fmt.Println("   Run 'sona config set api_key <YOUR_KEY>' to set it")
		}

		// Check output directory
		fmt.Println("\n4. Default Output Directory:")
		defaultPath := config.GetOutputPath()
		fmt.Printf("   %s\n", defaultPath)

		// Check if directory exists and is writable
		if info, err := os.Stat(defaultPath); err == nil && info.IsDir() {
			if testFile := os.WriteFile(filepath.Join(defaultPath, ".test"), []byte("test"), 0644); testFile == nil {
				os.Remove(filepath.Join(defaultPath, ".test"))
				fmt.Println("   Directory exists and is writable")
			} else {
				fmt.Println("   Directory exists but may not be writable")
			}
		} else {
			fmt.Println("   Directory does not exist (will be created automatically)")
		}

		fmt.Println("\nStatus check completed!")
	},
}

func main() {
	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.CloseLogger()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
