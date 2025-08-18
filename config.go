package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Manage configuration settings for the transcript converter tool.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. Available keys: api_key`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		switch key {
		case "api_key":
			viper.Set("assemblyai.api_key", value)
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}
			fmt.Printf("API key saved successfully!\n")
		default:
			fmt.Printf("Unknown config key: %s\n", key)
		}
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current Configuration:")
		fmt.Printf("API Key: %s\n", maskAPIKey(viper.GetString("assemblyai.api_key")))
		fmt.Printf("Config File: %s\n", viper.ConfigFileUsed())
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
}

func initConfig() {
	// Set default config file path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	configDir := filepath.Join(home, ".transcript-converter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		return
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("assemblyai.api_key", "")
	viper.SetDefault("output.default_path", filepath.Join(home, "transcripts"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: %v\n", err)
		}
	}

	// Check for environment variable
	if apiKey := os.Getenv("ASSEMBLYAI_API_KEY"); apiKey != "" {
		viper.Set("assemblyai.api_key", apiKey)
	}
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}

func getAPIKey() string {
	apiKey := viper.GetString("assemblyai.api_key")
	if apiKey == "" {
		fmt.Println("Error: AssemblyAI API key not found!")
		fmt.Println("Please set it using one of these methods:")
		fmt.Println("1. Set environment variable: export ASSEMBLYAI_API_KEY='your_key_here'")
		fmt.Println("2. Use config command: transcript-converter config set api_key 'your_key_here'")
		os.Exit(1)
	}
	return apiKey
}
