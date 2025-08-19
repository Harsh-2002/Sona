# sona - audio transcription tool

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

# Option B: Using the tool
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
├── docs/               # Documentation
├── examples/           # Usage examples
├── scripts/            # Build and utility scripts
├── build/              # Build outputs
├── go.mod              # Go module definition
└── README.md           # This file
```

### Package Details

- **`cmd/sona/`** - CLI application using Cobra framework
- **`pkg/assemblyai/`** - HTTP client for AssemblyAI REST API
- **`pkg/config/`** - Configuration management with Viper
- **`pkg/transcriber/`** - Orchestrates the transcription process
- **`pkg/youtube/`** - YouTube audio download using yt-dlp

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

# Run
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

The tool creates a config file at `~/.sona/config.toml`:

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
# or
sona config set api_key "your_key_here"
```

**YouTube Download Fails**
- Check internet connection
- Video may be private/restricted
- Ensure yt-dlp is installed or can be auto-installed

**Audio Format Issues**
- The tool will attempt to convert audio to MP3 format using FFmpeg
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

We use semantic versioning (SemVer) for releases. The project includes a version management script:

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
