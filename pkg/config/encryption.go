package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
)

// EncryptionManager handles encryption/decryption of sensitive config data
type EncryptionManager struct {
	masterKey []byte
}

// NewEncryptionManager creates a new encryption manager with a system-derived master key
func NewEncryptionManager() (*EncryptionManager, error) {
	// Generate a master key based on system information
	masterKey, err := generateMasterKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key: %v", err)
	}

	return &EncryptionManager{
		masterKey: masterKey,
	}, nil
}

// generateMasterKey creates a deterministic master key based on system information
func generateMasterKey() ([]byte, error) {
	// Get system information to create a unique but deterministic key
	systemInfo := fmt.Sprintf("%s-%s-%s-%s",
		runtime.GOOS,           // Operating system
		runtime.GOARCH,         // Architecture
		getHostname(),          // Hostname
		getUsername(),          // Username
	)

	// Create SHA256 hash of system info
	hash := sha256.Sum256([]byte(systemInfo))
	return hash[:], nil
}

// getHostname returns the system hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown-host"
	}
	return hostname
}

// getUsername returns the current username
func getUsername() string {
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	if username == "" {
		username = "unknown-user"
	}
	return username
}

// Encrypt encrypts a string value using AES-256-GCM
func (em *EncryptionManager) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Create AES cipher
	block, err := aes.NewCipher(em.masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts an encrypted string value
func (em *EncryptionManager) Decrypt(encryptedText string) (string, error) {
	if encryptedText == "" {
		return "", nil
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(em.masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a string appears to be encrypted
func (em *EncryptionManager) IsEncrypted(text string) bool {
	if text == "" {
		return false
	}
	
	// Try to decode as base64 and check if it's long enough to be encrypted
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return false
	}
	
	// Encrypted text should be at least 28 bytes (12 nonce + 16 tag + some data)
	return len(decoded) >= 28
}
