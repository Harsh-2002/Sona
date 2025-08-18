package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	// Load environment variables from .env file if it exists
	godotenv.Load()

	var rootCmd = &cobra.Command{
		Use:   "sona",
		Short: "Convert audio files and YouTube videos to text using AssemblyAI",
		Long: `A CLI tool that converts audio files and YouTube videos to text transcripts using AssemblyAI.
		
Features:
- Transcribe local audio files
- Download and transcribe YouTube videos
- Save transcripts to custom or default paths
- Progress tracking and nice CLI experience`,
	}

	// Add commands
	rootCmd.AddCommand(transcribeCmd)
	rootCmd.AddCommand(configCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
