# Sona - Audio Transcription Tool

A powerful, **completely independent** CLI tool that converts audio files and YouTube videos to text transcripts using AssemblyAI's advanced speech recognition technology.

## ✨ Features

- **Local Audio Files** - Transcribe any supported audio format
- **YouTube Videos** - Download and transcribe YouTube videos
- **Advanced AI** - Uses AssemblyAI's latest speech models
- **Smart Configuration** - API keys via environment or config commands
- **Flexible Output** - Auto-generate filenames or specify custom paths
- **Secure Storage** - API keys are stored with encryption
- **Interactive Mode** - Guided experience with step-by-step prompts and remembered settings

## 🚀 Quick Start

### 1. Download & Run
```bash
# Download the binary for your platform
# Then run immediately - no installation needed!
./sona --help
```

### 2. Set API Key
```bash
# Option A: Environment variable
export ASSEMBLYAI_API_KEY="your_api_key_here"

# Option B: Using Sona
./sona config set api_key "your_api_key_here"
```

### 3. Start Transcribing
```bash
# YouTube video
./sona transcribe "https://youtube.com/watch?v=dQw4w9WgXcQ"

# Local audio file
./sona transcribe "./audio.mp3"
```

## 📋 Requirements

- **AssemblyAI API Key** - [Get one here](https://www.assemblyai.com/)
- **yt-dlp** - For YouTube downloads (auto-installed if not found)
- **FFmpeg** - For audio format conversion (auto-installed if not found)

## 🏗️ Project Structure

```
sona-ai/
├── cmd/sona/           # Main application entry point
├── pkg/                # Core packages
│   ├── assemblyai/     # AssemblyAI API client
│   ├── config/         # Configuration management
│   ├── transcriber/    # Main transcription logic
│   └── youtube/        # YouTube download (pure Go)
├── scripts/            # Build and utility scripts
├── build/              # Build outputs
├── go.mod              # Go module definition
└── README.md           # This file
```

### Package Details

- **`cmd/sona/`** - CLI application using Cobra framework
- **`pkg/assemblyai/`** - HTTP client for AssemblyAI REST API
- **`pkg/config/** - Configuration management with Viper
- **`pkg/transcriber/`** - Orchestrates the transcription process
- **`pkg/youtube/`** - YouTube audio download using yt-dlp

## 🏛️ Architecture

### Architecture Principles

- **Separation of Concerns** - Each package has a single responsibility
- **Dependency Injection** - Dependencies are injected rather than hardcoded
- **Interface-based Design** - Uses interfaces for flexibility and testing
- **Error Handling** - Comprehensive error handling with meaningful messages
- **Configuration Management** - Centralized configuration with multiple sources

### Package Architecture

#### `pkg/assemblyai`
**Responsibility**: Communication with AssemblyAI REST API

**Key Components**:
- `Client` - HTTP client with authentication
- `TranscribeAudio()` - Main transcription method
- `uploadAudioFile()` - File upload to AssemblyAI
- `submitTranscription()` - Submit transcription request
- `pollTranscription()` - Poll for completion

**Dependencies**: Standard library (`net/http`, `encoding/json`, `mime/multipart`)

#### `pkg/config`
**Responsibility**: Configuration management and persistence

**Key Components**:
- `InitConfig()` - Initialize configuration system
- `GetAPIKey()` - Retrieve API key with validation
- `GetOutputPath()` - Get default output directory
- `ConfigCmd` - CLI commands for configuration

**Dependencies**: `github.com/spf13/viper`, `github.com/spf13/cobra`

#### `pkg/transcriber`
**Responsibility**: Orchestrating the transcription process

**Key Components**:
- `TranscribeCmd` - Main CLI command
- `processYouTubeVideo()` - Handle YouTube URLs
- `processLocalAudio()` - Handle local files
- `saveTranscript()` - Save results with smart naming

**Dependencies**: All other packages, `github.com/spf13/cobra`

#### `pkg/youtube`
**Responsibility**: YouTube audio download using pure Go

**Key Components**:
- `DownloadAudio()` - Download audio from YouTube URL
- `IsYouTubeURL()` - Validate YouTube URLs
- `GetVideoInfo()` - Get video metadata

**Dependencies**: Native yt-dlp binary management

### Data Flow

```
User Input → CLI Command → Transcriber → YouTube/AssemblyAI → Output File
     ↓              ↓           ↓              ↓              ↓
  YouTube URL   transcribe   Determine    Download &     Save to
  or File Path              Source Type   Transcribe     File System
```

### Configuration Flow

```
Environment Variable → Viper → Config File → Runtime
       ↓                ↓         ↓          ↓
  ASSEMBLYAI_API_KEY  Load    ~/.sona/   Use in App
```

### Error Handling Strategy

1. **Input Validation** - Validate user input early
2. **Graceful Degradation** - Provide helpful error messages
3. **Context Preservation** - Include relevant context in errors
4. **User Guidance** - Suggest solutions for common issues

### Security Considerations

- **API Key Management** - Secure storage and masking
- **File Validation** - Validate audio file types and sizes
- **URL Validation** - Sanitize YouTube URLs
- **Output Path Validation** - Prevent path traversal attacks

### Performance Considerations

- **Streaming Downloads** - Download audio in chunks
- **Concurrent Processing** - Process multiple files simultaneously
- **Memory Management** - Efficient handling of large audio files
- **Connection Pooling** - Reuse HTTP connections

## 🛠️ Building from Source

### Prerequisites
- Go 1.22 or later

### Build Steps
```bash
# Clone the repository
git clone <repository-url>
cd sona-ai

# Install dependencies
go mod tidy

# Build the binary
go build -o build/sona cmd/sona/main.go

# Run Sona
./build/sona --help
```

## 📖 Usage

### Basic Commands

```bash
# Interactive mode (default when no arguments provided)
sona

# Show help
sona --help

# Transcribe YouTube video
sona transcribe "https://youtube.com/watch?v=..."

# Transcribe local file
sona transcribe "./audio.mp3"

# With custom output
sona transcribe "video.mp4" --output transcript.txt

# With specific model
sona transcribe "audio.mp3" --model slam-1
```

## 🚀 Usage Examples

### Quick Start

#### 1. Set API Key
```bash
# Set your AssemblyAI API key
export ASSEMBLYAI_API_KEY="your_api_key_here"
```

#### 2. Transcribe YouTube Video
```bash
# Use Sona to transcribe YouTube video
./sona transcribe "https://youtube.com/watch?v=dQw4w9WgXcQ"
```

#### 3. Transcribe Local Audio File
```bash
# Use Sona to transcribe local audio
./sona transcribe "./meeting_recording.mp3"
```

### Advanced Usage

#### Custom Output Path
```bash
./sona transcribe "video.mp4" --output ./my_transcript.txt
```

#### Different Speech Model
```bash
# Use best quality model
./sona transcribe "audio.mp3" --model best

# Use fastest model
./sona transcribe "audio.mp3" --model nano
```

#### Configuration Management
```bash
# Show current config
./sona config show

# Set API key via config
./sona config set api_key "your_new_key_here"
```

### Output Examples

#### YouTube Video Output
```
🎥 Detected YouTube URL, downloading audio...
📥 Downloading audio from YouTube using pure Go...
🎬 Video: Rick Astley - Never Gonna Give You Up
⏱️  Duration: 3m 33s
🎵 Audio format: AUDIO_QUALITY_MEDIUM
⬇️  Downloading audio stream...
✅ Audio download completed successfully!
✅ Downloaded audio to: /tmp/sona-12345/audio.mp4
🔊 Starting transcription with AssemblyAI...
📤 Audio file uploaded successfully
📝 Transcription submitted, waiting for completion...
⏳ Status: queued, waiting...
⏳ Status: processing, waiting...
✅ Transcription completed successfully!
✅ Transcript saved to: /home/user/transcripts/youtube_dQw4w9WgXcQ.txt
📝 Transcript length: 1234 characters
```

#### Local File Output
```
🎵 Processing local audio file...
🔊 Starting transcription with AssemblyAI...
📤 Audio file uploaded successfully
📝 Transcription submitted, waiting for completion...
⏳ Status: queued, waiting...
⏳ Status: processing, waiting...
✅ Transcription completed successfully!
✅ Transcript saved to: /home/user/transcripts/meeting_recording_transcript.txt
📝 Transcript length: 5678 characters
```

### Configuration

```bash
# Show current config
sona config show

# Set API key
sona config set api_key "your_key_here"
```

### Command Options

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Custom output file path | Auto-generated |
| `-m, --model` | Speech model (slam-1, best, nano) | slam-1 |

## 🎯 Speech Models

- **`slam-1`** (default) - Latest model, best accuracy
- **`best`** - High accuracy, good for most use cases  
- **`nano`** - Fastest, good for real-time applications

## 📁 Output

By default, transcripts are saved to:
- `~/sona/[title]-[date].txt`

Where:
- `title` is a simplified version of the YouTube video title or local file name
- `date` is the current date (format: YYYYMMDD)

Override with `--output` flag.

## 🔧 Configuration

Sona creates a config file at `~/.sona/config.toml`:

```toml
[assemblyai]
api_key = "encrypted_api_key_here"

[output]
default_path = "/home/user/sona"

[last_session]
source_type = "youtube"
speech_model = "nano"
output_path = "/custom/path/transcript.txt"
```

## 🔒 Security

API keys are automatically encrypted using AES-256-GCM encryption with a system-derived master key. This ensures your API keys are securely stored and can only be decrypted on the same system they were encrypted on.

## 🐛 Troubleshooting

### Common Issues

**API Key Not Found**
```bash
# Set your API key
export ASSEMBLYAI_API_KEY="your_key_here"
# or use Sona's config command
sona config set api_key "your_key_here"
```

**YouTube Download Fails**
- Check internet connection
- Video may be private/restricted
- Sona will auto-install yt-dlp if not found

**Audio Format Issues**
- Sona will attempt to convert audio to MP3 format using FFmpeg
- If conversion fails, try converting the file manually to MP3 format
- Some DRM-protected or specialized formats may not be convertible

**Transcription Fails**
- Verify API key is correct
- Check audio file format and size
- Ensure file is not corrupted

### Supported Audio Formats
- MP3, WAV, M4A, FLAC, OGG
- Video files with audio (MP4, AVI, MOV)
- Max file size: 1GB
- Min duration: 160ms

## 🚀 Version Management

Sona uses semantic versioning (SemVer) for releases. The project includes a version management script:

```bash
# Show current version
./scripts/version.sh show

# Show next version options
./scripts/version.sh next

# Create new version (automatically triggers GitHub Actions release)
./scripts/version.sh create patch    # 1.0.0 → 1.0.1
./scripts/version.sh create minor    # 1.0.0 → 1.1.0
./scripts/version.sh create major    # 1.0.0 → 2.0.0
```

### Release Process
1. **Create version**: `./scripts/version.sh create patch`
2. **Script automatically**: Creates git tag and pushes to remote
3. **GitHub Actions**: Automatically builds and releases for all platforms
4. **Result**: New release with binaries for Linux, macOS, and Windows

## 🚀 Future Extensibility

### Planned Features
- **Batch Processing** - Process multiple files
- **Audio Preprocessing** - Audio enhancement before transcription
- **Multiple Output Formats** - JSON, SRT, VTT
- **Progress Bars** - Visual download/transcription progress
- **Resume Support** - Resume interrupted downloads
- **Auto-Update** - Self-updating binary from GitHub releases

### Architecture Considerations
- **Plugin System** - Support for custom audio sources
- **Queue System** - Background processing of multiple files
- **Caching** - Cache downloaded audio files
- **Metrics** - Performance and usage metrics

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- [AssemblyAI](https://www.assemblyai.com/) - Speech recognition API
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - YouTube download tool
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management

## 📞 Support

- Check the troubleshooting section above
- Search existing GitHub issues
- Create a new issue with detailed information
