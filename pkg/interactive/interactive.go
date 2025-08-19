package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Harsh-2002/Sona/pkg/config"
	"github.com/Harsh-2002/Sona/pkg/transcriber"
	"github.com/Harsh-2002/Sona/pkg/youtube"
	"github.com/spf13/cobra"
)

// InteractiveCmd represents the interactive command
var InteractiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive mode",
	Long:  `Start interactive mode to guide you through the transcription process step by step.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveMode()
	},
}

func runInteractiveMode() {
	fmt.Println("--------------------------------")
	fmt.Println("❇️  Sona is your go-to tool for turning audio files or YouTube videos into text—fast, easy, and accurate.")
	fmt.Println("--------------------------------")

	// Check if API key is set
	apiKey := checkAndSetAPIKey()

	// Get last session settings
	lastSourceType := config.GetLastSourceType()
	lastSpeechModel := config.GetLastSpeechModel()
	lastOutputPath := config.GetLastOutputPath()

	// Ask for source type with last used as default
	sourceType := promptSourceType(lastSourceType)

	// Get source path or URL
	source := promptSource(sourceType)

	// Ask for output path (optional) with last used as default
	outputPath := promptOutputPath(lastOutputPath)

	// Ask for speech model (optional) with last used as default
	speechModel := promptSpeechModel(lastSpeechModel)

	// Confirm settings
	if !confirmSettings(sourceType, source, outputPath, speechModel) {
		fmt.Println("Operation canceled")
		return
	}

	// Save settings for next time
	config.SaveLastSession(sourceType, speechModel, outputPath)

	// Set command-line flags
	if outputPath != "" {
		transcriber.SetOutputPath(outputPath)
	}
	if speechModel != "" {
		transcriber.SetSpeechModel(speechModel)
	}

	// Process based on source type
	var err error
	if sourceType == "youtube" {
		err = transcriber.ProcessYouTubeVideo(source, apiKey)
	} else {
		err = transcriber.ProcessLocalAudio(source, apiKey)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

// checkAndSetAPIKey checks if API key is set and prompts user to set it if not
func checkAndSetAPIKey() string {
	apiKey := ""

	// Try to get existing API key
	apiKey = config.GetAPIKeyNoExit()

	// If no API key, prompt user to enter one
	if apiKey == "" {
		fmt.Println("\nNo AssemblyAI API key found. You need an API key to use this tool.")
		fmt.Println("You can get one for free at https://www.assemblyai.com/")

		for {
			fmt.Print("\nPlease enter your AssemblyAI API key: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			apiKey = strings.TrimSpace(scanner.Text())

			if apiKey == "" {
				fmt.Println("API key cannot be empty. Please try again.")
				continue
			}

			// Save the API key
			fmt.Print("Do you want to save this API key for future use? (y/n): ")
			scanner.Scan()
			if strings.ToLower(strings.TrimSpace(scanner.Text())) == "y" {
				config.SaveAPIKey(apiKey)
				fmt.Println("API key saved successfully")
			}

			break
		}
	}

	return apiKey
}

// promptSourceType asks user to select source type
func promptSourceType(lastSourceType string) string {
	fmt.Println("\nWhat type of source would you like to transcribe?")
	fmt.Println("1. YouTube video")
	fmt.Println("2. Local audio file")

	// Show last used option if available
	defaultOption := ""
	if lastSourceType == "youtube" {
		defaultOption = "1"
		fmt.Println("Last used: YouTube video")
	} else if lastSourceType == "local" {
		defaultOption = "2"
		fmt.Println("Last used: Local audio file")
	}

	for {
		if defaultOption != "" {
			fmt.Printf("\nEnter your choice (1 or 2, press Enter for last used [%s]): ", defaultOption)
		} else {
			fmt.Print("\nEnter your choice (1 or 2): ")
		}

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		// Use default if empty
		if choice == "" && defaultOption != "" {
			choice = defaultOption
		}

		if choice == "1" {
			return "youtube"
		} else if choice == "2" {
			return "local"
		} else {
			fmt.Println("Invalid choice. Please enter 1 or 2.")
		}
	}
}

// promptSource asks user for source path or URL
func promptSource(sourceType string) string {
	var prompt string
	if sourceType == "youtube" {
		prompt = "Enter YouTube URL: "
	} else {
		prompt = "Enter path to audio file: "
	}

	for {
		fmt.Print("\n" + prompt)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		source := strings.TrimSpace(scanner.Text())

		if source == "" {
			fmt.Println("Source cannot be empty. Please try again.")
			continue
		}

		// Validate source
		if sourceType == "youtube" && !youtube.IsYouTubeURL(source) {
			fmt.Println("Invalid YouTube URL. Please enter a valid URL.")
			continue
		} else if sourceType == "local" {
			if _, err := os.Stat(source); os.IsNotExist(err) {
				fmt.Println("File not found. Please enter a valid path.")
				continue
			}
		}

		return source
	}
}

// promptOutputPath asks user for output path (optional)
func promptOutputPath(lastOutputPath string) string {
	prompt := "\nEnter output path (leave blank for default)"

	// Show last used path if available
	if lastOutputPath != "" {
		prompt += fmt.Sprintf(" or press Enter for last used [%s]", lastOutputPath)
	}

	fmt.Print(prompt + ": ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	path := strings.TrimSpace(scanner.Text())

	// Use last path if input is empty and last path exists
	if path == "" && lastOutputPath != "" {
		return lastOutputPath
	}

	return path
}

// promptSpeechModel asks user for speech model (optional)
func promptSpeechModel(lastModel string) string {
	fmt.Println("\nSelect speech model:")
	fmt.Println("1. slam-1 (best accuracy)")
	fmt.Println("2. best (good for most use cases)")
	fmt.Println("3. nano (fastest, good for real-time)")

	// Determine default choice based on last used model
	defaultChoice := ""
	defaultModel := "slam-1"

	switch lastModel {
	case "slam-1":
		defaultChoice = "1"
		defaultModel = "slam-1"
	case "best":
		defaultChoice = "2"
		defaultModel = "best"
	case "nano":
		defaultChoice = "3"
		defaultModel = "nano"
	}

	// Show last used model if available
	if defaultChoice != "" {
		fmt.Printf("Last used: %s\n", lastModel)
		fmt.Printf("\nEnter your choice (1-3, or press Enter for last used [%s]): ", defaultChoice)
	} else {
		fmt.Print("\nEnter your choice (1-3, or leave blank for default): ")
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	// Use default if empty
	if choice == "" {
		if defaultChoice != "" {
			return defaultModel
		}
		return "slam-1"
	}

	switch choice {
	case "1":
		return "slam-1"
	case "2":
		return "best"
	case "3":
		return "nano"
	default:
		fmt.Println("Invalid choice. Using default model (slam-1).")
		return "slam-1"
	}
}

// confirmSettings shows a summary and asks user to confirm
func confirmSettings(sourceType, source, outputPath, speechModel string) bool {
	fmt.Println("\nSummary of settings:")
	fmt.Printf("Source type: %s\n", sourceType)
	fmt.Printf("Source: %s\n", source)

	if outputPath != "" {
		fmt.Printf("Output path: %s\n", outputPath)
	} else {
		fmt.Println("Output path: [default]")
	}

	fmt.Printf("Speech model: %s\n", speechModel)

	fmt.Print("\nProceed with these settings? (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ToLower(strings.TrimSpace(scanner.Text())) == "y"
}
