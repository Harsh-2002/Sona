package assemblyai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type TranscriptionRequest struct {
	AudioURL    string `json:"audio_url"`
	SpeechModel string `json:"speech_model"`
}

type TranscriptionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type TranscriptResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Text   string `json:"text"`
	Error  string `json:"error,omitempty"`
}

// Client represents an AssemblyAI client
type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new AssemblyAI client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// TranscribeAudio transcribes an audio file using AssemblyAI
func (c *Client) TranscribeAudio(audioPath string, speechModel string) (string, error) {
	fmt.Println("🔊 Starting transcription with AssemblyAI...")

	// First, upload the audio file
	uploadURL, err := c.uploadAudioFile(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to upload audio file: %v", err)
	}

	fmt.Println("📤 Audio file uploaded successfully")

	// Submit transcription request
	transcriptID, err := c.submitTranscription(uploadURL, speechModel)
	if err != nil {
		return "", fmt.Errorf("failed to submit transcription: %v", err)
	}

	fmt.Println("📝 Transcription submitted, waiting for completion...")

	// Poll for completion
	transcript, err := c.pollTranscription(transcriptID)
	if err != nil {
		return "", fmt.Errorf("failed to get transcription: %v", err)
	}

	if transcript.Status == "error" {
		return "", fmt.Errorf("transcription failed: %s", transcript.Error)
	}

	fmt.Println("✅ Transcription completed successfully!")
	return transcript.Text, nil
}

// uploadAudioFile uploads an audio file to AssemblyAI and returns the upload URL
func (c *Client) uploadAudioFile(audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %v", err)
	}
	defer file.Close()

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	part, err := writer.CreateFormFile("file", "audio.mp3")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file data: %v", err)
	}

	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make upload request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var uploadResp struct {
		UploadURL string `json:"upload_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("failed to decode upload response: %v", err)
	}

	return uploadResp.UploadURL, nil
}

// submitTranscription submits a transcription request to AssemblyAI
func (c *Client) submitTranscription(audioURL string, speechModel string) (string, error) {
	request := TranscriptionRequest{
		AudioURL:    audioURL,
		SpeechModel: speechModel,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to submit transcription: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("transcription submission failed with status %d: %s", resp.StatusCode, string(body))
	}

	var transcriptResp TranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&transcriptResp); err != nil {
		return "", fmt.Errorf("failed to decode transcription response: %v", err)
	}

	return transcriptResp.ID, nil
}

// pollTranscription polls the transcription status until completion
func (c *Client) pollTranscription(transcriptID string) (*TranscriptResult, error) {
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create polling request: %v", err)
		}

		req.Header.Set("Authorization", c.APIKey)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to poll transcription: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("polling failed with status %d", resp.StatusCode)
		}

		var result TranscriptResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode polling response: %v", err)
		}
		resp.Body.Close()

		switch result.Status {
		case "completed":
			return &result, nil
		case "error":
			return &result, nil
		case "queued", "processing":
			fmt.Printf("⏳ Status: %s, waiting...\n", result.Status)
			time.Sleep(3 * time.Second)
		default:
			fmt.Printf("⏳ Unknown status: %s, waiting...\n", result.Status)
			time.Sleep(3 * time.Second)
		}
	}
}
