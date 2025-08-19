package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/root/sona-ai/pkg/config"
	"github.com/root/sona-ai/pkg/interactive"
	"github.com/root/sona-ai/pkg/transcriber"
	"github.com/spf13/cobra"
)

// Version will be set by the build process
var version = "dev"

func main() {
	// Load environment variables from .env file if it exists
	godotenv.Load()

	// Initialize configuration
	config.InitConfig()

	var rootCmd = &cobra.Command{
		Use:     "sona",
		Version: version,
		Short:   "Audio transcription tool",
		Long: `Sona - Audio Transcription Tool

A CLI tool that converts audio files and YouTube videos to text transcripts using AssemblyAI.

Features:
- Transcribe local audio files
- Download and transcribe YouTube videos
- Save transcripts to custom or default paths
- Interactive mode for guided experience`,
		// If no subcommand is provided, run interactive mode
		Run: func(cmd *cobra.Command, args []string) {
			interactive.InteractiveCmd.Run(cmd, args)
		},
	}

	// Add commands
	rootCmd.AddCommand(transcriber.TranscribeCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(interactive.InteractiveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
