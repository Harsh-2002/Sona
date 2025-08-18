# Sona - Audio to Text Transcript Converter

A powerful, **completely independent** CLI tool that converts audio files and YouTube videos to text transcripts using AssemblyAI's advanced speech recognition technology.

## ✨ Features

- 🎵 **Local Audio Files** - Transcribe any supported audio format
- 🎥 **YouTube Videos** - Download and transcribe using pure Go (no external tools)
- 🤖 **Advanced AI** - Uses AssemblyAI's latest speech models
- ⚙️ **Smart Configuration** - API keys via environment or config commands
- 📁 **Flexible Output** - Auto-generate filenames or specify custom paths
- 🔒 **100% Independent** - Single binary, no external dependencies

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
- **No external tools needed** - Everything is included in the binary

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
- **`pkg/youtube/`** - Pure Go YouTube audio download using kkdai/youtube/v2

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
- **YouTube**: `~/transcripts/youtube_[VIDEO_ID].txt`
- **Local files**: `~/transcripts/[FILENAME]_transcript.txt`

Override with `--output` flag.

## 🔧 Configuration

The tool creates a config file at `~/.sona/config.toml`:

```toml
[assemblyai]
api_key = "your_api_key_here"

[output]
default_path = "/home/user/transcripts"
```

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
- Tool automatically handles YouTube's current systems

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
- [kkdai/youtube/v2](https://github.com/kkdai/youtube/v2) - Pure Go YouTube library
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management

## 📞 Support

- Check the troubleshooting section above
- Search existing GitHub issues
- Create a new issue with detailed information
