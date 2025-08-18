package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var encryptionManager *EncryptionManager
var configFilePath string

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Manage configuration settings for the sona tool.`,
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
			// Encrypt the API key if encryption is available
			if encryptionManager != nil {
				encryptedValue, err := encryptionManager.Encrypt(value)
				if err != nil {
					fmt.Printf("Warning: Could not encrypt API key: %v\n", err)
					fmt.Printf("API key will be stored in plain text\n")
					viper.Set("assemblyai.api_key", value)
				} else {
					viper.Set("assemblyai.api_key", encryptedValue)
					fmt.Printf("🔒 API key encrypted and saved successfully!\n")
				}
			} else {
				viper.Set("assemblyai.api_key", value)
				fmt.Printf("⚠️  API key saved in plain text (encryption not available)\n")
			}
			
			// Persist config: always write to ~/.sona/config.toml
			var err error
			if _, statErr := os.Stat(configFilePath); os.IsNotExist(statErr) {
				err = viper.WriteConfigAs(configFilePath)
			} else {
				err = viper.WriteConfig()
			}
			
			if err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}
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
		fmt.Printf("API Key: %s\n", MaskAPIKey(viper.GetString("assemblyai.api_key")))
		fmt.Printf("Config File: %s\n", viper.ConfigFileUsed())
	},
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configShowCmd)
}

// InitConfig initializes the configuration system
func InitConfig() {
	// Initialize encryption manager
	var err error
	encryptionManager, err = NewEncryptionManager()
	if err != nil {
		fmt.Printf("Warning: Could not initialize encryption: %v\n", err)
		fmt.Printf("API keys will be stored in plain text\n")
	}

	// Set default config file path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	configDir := filepath.Join(home, ".sona")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		return
	}

	configFilePath = filepath.Join(configDir, "config.toml")
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("toml")

	// Set defaults
	viper.SetDefault("assemblyai.api_key", "")
	viper.SetDefault("output.default_path", filepath.Join(home, "sona"))
	viper.SetDefault("last_session.source_type", "")
	viper.SetDefault("last_session.speech_model", "slam-1")
	viper.SetDefault("last_session.output_path", "")

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: %v\n", err)
		}
	}

	// Write default config if it doesn't exist
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := viper.WriteConfigAs(configFilePath); err != nil {
			fmt.Printf("Warning: Could not write default config file: %v\n", err)
		}
	}

	// Check for environment variable
	if apiKey := os.Getenv("ASSEMBLYAI_API_KEY"); apiKey != "" {
		viper.Set("assemblyai.api_key", apiKey)
	}
}

func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}

// GetAPIKey returns the AssemblyAI API key and exits if not found
func GetAPIKey() string {
	apiKey := GetAPIKeyNoExit()
	if apiKey == "" {
		fmt.Println("Error: AssemblyAI API key not found!")
		fmt.Println("Please set it using one of these methods:")
		fmt.Println("1. Set environment variable: export ASSEMBLYAI_API_KEY='your_key_here'")
		fmt.Println("2. Use config command: sona config set api_key 'your_key_here'")
		fmt.Println("3. Run in interactive mode: sona")
		os.Exit(1)
	}
	return apiKey
}

// GetAPIKeyNoExit returns the AssemblyAI API key without exiting if not found
func GetAPIKeyNoExit() string {
	apiKey := viper.GetString("assemblyai.api_key")
	
	// Check if API key is empty
	if apiKey == "" {
		return ""
	}

	// Decrypt the API key if it's encrypted
	if encryptionManager != nil && encryptionManager.IsEncrypted(apiKey) {
		decryptedKey, err := encryptionManager.Decrypt(apiKey)
		if err != nil {
			fmt.Printf("Error: Failed to decrypt API key: %v\n", err)
			fmt.Println("Please reset your API key using: sona config set api_key 'your_key_here'")
			return ""
		}
		return decryptedKey
	}

	return apiKey
}

// SaveAPIKey saves the API key to the config file
func SaveAPIKey(apiKey string) error {
	// Encrypt the API key if encryption is available
	if encryptionManager != nil {
		encryptedValue, err := encryptionManager.Encrypt(apiKey)
		if err != nil {
			fmt.Printf("Warning: Could not encrypt API key: %v\n", err)
			fmt.Printf("API key will be stored in plain text\n")
			viper.Set("assemblyai.api_key", apiKey)
		} else {
			viper.Set("assemblyai.api_key", encryptedValue)
		}
	} else {
		viper.Set("assemblyai.api_key", apiKey)
		fmt.Printf("Warning: API key saved in plain text (encryption not available)\n")
	}
	
	// Persist config
	var err error
	if _, statErr := os.Stat(configFilePath); os.IsNotExist(statErr) {
		err = viper.WriteConfigAs(configFilePath)
	} else {
		err = viper.WriteConfig()
	}
	
	return err
}

// GetOutputPath returns the default output path
func GetOutputPath() string {
	return viper.GetString("output.default_path")
}

// GetLastSourceType returns the last used source type
func GetLastSourceType() string {
	return viper.GetString("last_session.source_type")
}

// GetLastSpeechModel returns the last used speech model
func GetLastSpeechModel() string {
	model := viper.GetString("last_session.speech_model")
	if model == "" {
		return "slam-1" // Default if not set
	}
	return model
}

// GetLastOutputPath returns the last used output path
func GetLastOutputPath() string {
	return viper.GetString("last_session.output_path")
}

// SaveLastSession saves the last session settings
func SaveLastSession(sourceType, speechModel, outputPath string) error {
	viper.Set("last_session.source_type", sourceType)
	viper.Set("last_session.speech_model", speechModel)
	viper.Set("last_session.output_path", outputPath)
	
	// Persist config
	return viper.WriteConfig()
}
