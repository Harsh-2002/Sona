package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies for the current platform",
	Long:  "Install yt-dlp and FFmpeg dependencies for the current platform. This command will download and install the appropriate binaries for your operating system.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sona Dependency Installation")
		fmt.Println("============================")

		// Install yt-dlp
		fmt.Println("\n1. YouTube Download (yt-dlp):")
		fmt.Println("   Installing...")
		if err := youtube.InstallYtDlp(); err != nil {
			fmt.Printf("   Failed: %v\n", err)
			fmt.Println("   💡 Check logs at:", logger.GetLogPath())
			os.Exit(1)
		}
		fmt.Println("   ✅ Installed successfully")

		// Install FFmpeg
		fmt.Println("\n2. Audio Processing (FFmpeg):")
		fmt.Println("   Installing...")
		if err := transcriber.InstallFFmpeg(); err != nil {
			fmt.Printf("   Failed: %v\n", err)
			fmt.Println("   💡 Check logs at:", logger.GetLogPath())
			os.Exit(1)
		}
		fmt.Println("   ✅ Installed successfully")

		// On macOS, also check for ffprobe
		if runtime.GOOS == "darwin" {
			fmt.Println("\n3. macOS Audio Tools (ffprobe):")
			if _, err := transcriber.FindBinary("ffprobe"); err != nil {
				fmt.Println("   ⚠️  ffprobe not found after FFmpeg installation")
				fmt.Println("   💡 This might cause issues with YouTube downloads")
			} else {
				fmt.Println("   ✅ Available")
			}
		}

		fmt.Println("\nInstallation completed!")
		fmt.Println("💡 Run 'sona status' to verify the installation")
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
	rootCmd.AddCommand(installCmd)
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
			fmt.Println("   Not found (run 'sona install' to install)")
		}

		// Check FFmpeg
		fmt.Println("\n2. Audio Processing (FFmpeg):")
		if ffmpegPath, err := transcriber.FindBinary("ffmpeg"); err == nil {
			fmt.Printf("   FFmpeg available at: %s\n", ffmpegPath)

			// On macOS, also check for ffprobe
			if runtime.GOOS == "darwin" {
				if ffprobePath, err := transcriber.FindBinary("ffprobe"); err == nil {
					fmt.Printf("   ffprobe available at: %s\n", ffprobePath)
				} else {
					fmt.Println("   ffprobe not found (run 'sona install' to install)")
				}
			}
		} else {
			fmt.Println("   Not found (run 'sona install' to install)")
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
