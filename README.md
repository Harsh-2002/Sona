# Sona - Audio to Text Transcript Converter

A powerful CLI tool that converts audio files and YouTube videos to text transcripts using AssemblyAI's advanced speech recognition technology.

## Features

- 🎵 **Local Audio Files**: Transcribe any supported audio format (MP3, WAV, M4A, etc.)
- 🎥 **YouTube Videos**: Download and transcribe YouTube videos automatically
- 🤖 **Advanced AI**: Uses AssemblyAI's latest speech models for accurate transcription
- ⚙️ **Flexible Configuration**: Set API keys via environment variables or config commands
- 📁 **Smart Output**: Auto-generate filenames or specify custom output paths
- 🚀 **Fast & Efficient**: Optimized YouTube audio extraction and streaming transcription

## Features

- 🎵 **Local Audio Files**: Transcribe any supported audio format (MP3, WAV, M4A, etc.)
- 🎥 **YouTube Videos**: Download and transcribe YouTube videos automatically using pure Go
- 🤖 **Advanced AI**: Uses AssemblyAI's latest speech models for accurate transcription
- ⚙️ **Flexible Configuration**: Set API keys via environment variables or config commands
- 📁 **Smart Output**: Auto-generate filenames or specify custom output paths
- 🚀 **Fast & Efficient**: Pure Go implementation with no external dependencies
- 🔒 **Completely Independent**: Single binary that works anywhere - no external tools needed!

## No External Dependencies Required!

This tool is **100% independent** and **self-contained**:
- ✅ **No yt-dlp installation needed**
- ✅ **No youtube-dl installation needed** 
- ✅ **No external tools required**
- ✅ **Works on any platform** (Windows, macOS, Linux)
- ✅ **Single binary** - download and run immediately

## Installation

### Option 1: Build from Source (Recommended)

1. **Clone the repository:**
```bash
git clone <repository-url>
cd sona
```

2. **Install Go dependencies:**
```bash
go mod tidy
```

3. **Build the binary:**
```bash
go build -o sona
```

4. **Install globally (optional):**
```bash
sudo mv sona /usr/local/bin/
```

### Option 2: Download Pre-built Binary

Download the latest release binary for your platform from the releases page.

## Setup

### 1. Get AssemblyAI API Key

1. Sign up at [AssemblyAI](https://www.assemblyai.com/)
2. Go to your [API Keys page](https://www.assemblyai.com/app/api-keys)
3. Copy your API key

### 2. Configure the API Key

**Option A: Environment Variable (Recommended)**
```bash
export ASSEMBLYAI_API_KEY="your_api_key_here"
```

**Option B: Using the CLI config command**
```bash
./sona config set api_key "your_api_key_here"
```

**Option C: .env file**
Create a `.env` file in your project directory:
```env
ASSEMBLYAI_API_KEY=your_api_key_here
```

## Usage

### Basic Commands

**Transcribe a YouTube video:**
```bash
sona transcribe "https://youtube.com/watch?v=dQw4w9WgXcQ"
```

**Transcribe a local audio file:**
```bash
sona transcribe "./audio.mp3"
```

**Transcribe with custom output path:**
```bash
sona transcribe "https://youtube.com/watch?v=..." --output ./my_transcript.txt
```

**Transcribe with specific speech model:**
```bash
sona transcribe "./audio.mp3" --model slam-1
```

### Configuration Commands

**Show current configuration:**
```bash
sona config show
```

**Set configuration values:**
```bash
sona config set api_key "your_new_api_key"
```

### Command Options

- `--output, -o`: Specify custom output file path
- `--model, -m`: Choose speech model (slam-1, best, nano)
- `--help, -h`: Show help information

## Speech Models

- **slam-1** (default): Latest prompt-based model, best accuracy
- **best**: High accuracy, good for most use cases
- **nano**: Fastest, good for real-time applications

## Output

By default, transcripts are saved to:
- **YouTube videos**: `~/transcripts/youtube_[VIDEO_ID].txt`
- **Local files**: `~/transcripts/[FILENAME]_transcript.txt`

You can override this with the `--output` flag.

## Examples

### Example 1: Transcribe a YouTube Video
```bash
sona transcribe "https://youtube.com/watch?v=dQw4w9WgXcQ"
```

Output:
```
🎥 Detected YouTube URL, downloading audio...
📥 Downloading audio from YouTube...
✅ Downloaded audio to: /tmp/transcript-converter-12345/audio.mp3
🔊 Starting transcription with AssemblyAI...
📤 Audio file uploaded successfully
📝 Transcription submitted, waiting for completion...
⏳ Status: queued, waiting...
⏳ Status: processing, waiting...
✅ Transcription completed successfully!
✅ Transcript saved to: /home/user/transcripts/youtube_dQw4w9WgXcQ.txt
📝 Transcript length: 1234 characters
```

### Example 2: Transcribe Local Audio with Custom Output
```bash
sona transcribe "./meeting_recording.mp3" --output ./meeting_summary.txt
```

### Example 3: Use Different Speech Model
```bash
sona transcribe "./podcast.mp3" --model best
```

## Troubleshooting

### Common Issues

**1. "AssemblyAI API key not found" error:**
- Set your API key using one of the configuration methods above
- Check that the environment variable is properly exported

**2. YouTube download fails:**
- Check your internet connection
- Some videos may have restrictions or be private
- The tool automatically handles YouTube's current systems

**4. Transcription fails:**
- Verify your API key is correct
- Check that the audio file is in a supported format
- Ensure the audio file is not corrupted

### Supported Audio Formats

- MP3, WAV, M4A, FLAC, OGG
- Video files with audio (MP4, AVI, MOV)
- Maximum file size: 1GB
- Minimum duration: 160ms

## Configuration File

The tool creates a configuration file at `~/.transcript-converter/config.yaml`:

```yaml
assemblyai:
  api_key: "your_api_key_here"
output:
  default_path: "/home/user/transcripts"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

If you encounter any issues or have questions:

1. Check the troubleshooting section above
2. Search existing GitHub issues
3. Create a new issue with detailed information

## Acknowledgments

- [AssemblyAI](https://www.assemblyai.com/) for providing the speech recognition API
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) for efficient YouTube audio extraction
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
