# Sona - Audio Transcription Tool (CLI)

Turn audio files and YouTube videos into text with just a few commands. Sona uses AssemblyAI's speech recognition to give you accurate transcripts quickly and easily.

## ✨ What Sona Does

- **Audio Files** - Convert any audio file (MP3, WAV, M4A, etc.) to text
- **YouTube Videos** - Download and transcribe YouTube videos automatically
- **Smart AI** - Uses the latest speech recognition models for best accuracy
- **Easy Setup** - Simple configuration with your API key
- **Flexible Output** - Save transcripts wherever you want with smart naming
- **Secure** - Your API keys are encrypted and safe
- **Interactive** - Step-by-step guidance when you need it

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

## 📋 What You Need

- **AssemblyAI API Key** - [Get a free one here](https://www.assemblyai.com/)
- **Nothing else!** - Sona automatically installs everything else it needs

## 🏗️ How Sona Works

Sona is built with a clean, modular design that makes it easy to maintain and extend. The main components work together to handle everything from YouTube downloads to AI transcription.

## 🏛️ How Sona is Built

Sona is designed to be simple, reliable, and easy to use. Here's how it works under the hood:

### Main Components

- **AssemblyAI Client** - Handles communication with the speech recognition service
- **Configuration Manager** - Stores your API key securely and remembers your preferences
- **Transcription Engine** - Coordinates the entire process from start to finish
- **YouTube Downloader** - Automatically downloads and processes YouTube videos

### How It Works

1. **You run a command** - Tell Sona what to transcribe
2. **Sona figures out what to do** - YouTube video or local file?
3. **Downloads/processes audio** - Gets the audio ready for transcription
4. **Sends to AssemblyAI** - Uses their AI to convert speech to text
5. **Saves your transcript** - Puts the text where you want it

### Smart Features

- **Auto-installation** - Sona installs missing tools automatically
- **Error handling** - Clear messages when something goes wrong
- **Secure storage** - Your API keys are encrypted and safe
- **Memory management** - Handles large files efficiently

## 🛠️ Building from Source

If you want to build Sona yourself, you'll need Go 1.22 or later.

```bash
# Get the code
git clone <repository-url>
cd sona-ai

# Build Sona
go build -o build/sona cmd/sona/main.go

# Try it out
./build/sona --help
```

## 📖 How to Use Sona

```bash
# Get help
sona --help

# Transcribe a YouTube video
sona transcribe "https://youtube.com/watch?v=..."

# Transcribe an audio file
sona transcribe "./audio.mp3"

# Save to a specific location
sona transcribe "video.mp4" --output transcript.txt

# Use a different AI model
sona transcribe "audio.mp3" --model best
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

## 🎯 AI Models

Sona uses different AI models depending on your needs:

- **`slam-1`** (default) - Best accuracy, great for important content
- **`best`** - High accuracy, good for most things
- **`nano`** - Fastest, good when you need speed

## 📁 Where Transcripts Go

Sona automatically saves your transcripts to `~/sona/` with smart names like:
- `video-title-20241201.txt` for YouTube videos
- `filename-20241201.txt` for audio files

Use `--output` to save somewhere else.

## 🔧 Settings

Sona remembers your preferences in `~/.sona/config.toml`. You can:

- Set your API key once and forget about it
- Choose where to save transcripts by default
- Remember your last used settings

## 🔒 Keeping Your Data Safe

Sona encrypts your API keys automatically, so they're safe even if someone gets access to your computer. The encryption is tied to your specific system.

## 🐛 When Things Go Wrong

### Common Problems

**"API key not found"**
```bash
# Set your API key
export ASSEMBLYAI_API_KEY="your_key_here"
# or use Sona's config
sona config set api_key "your_key_here"
```

**YouTube won't download**
- Check your internet connection
- Make sure the video isn't private
- Sona will install the tools it needs automatically

**Audio won't convert**
- Sona tries to convert most formats to MP3
- If it fails, try converting the file yourself first
- Some protected files can't be converted

**Transcription fails**
- Check your API key is correct
- Make sure the audio file isn't corrupted
- Try a smaller file if it's very large

### What Sona Can Handle
- **Audio**: MP3, WAV, M4A, FLAC, OGG
- **Video**: MP4, AVI, MOV (audio will be extracted)
- **Size**: Up to 1GB
- **Length**: At least 160ms

## 🚀 What's Coming Next

We're planning to add:

- **Batch processing** - Handle multiple files at once
- **More output formats** - JSON, SRT, VTT files
- **Progress bars** - See how things are going
- **Auto-updates** - Sona updates itself automatically

## 📦 Automated Releases

When we create a new version, GitHub automatically:
- Builds Sona for all platforms (Linux, macOS, Windows)
- Creates a new release with all binaries
- Generates checksums for security
- Makes it easy to download the right version for your system

### 🚀 Creating Releases

Use our simple release script:
```bash
# Patch release (1.0.0 → 1.0.1)
./scripts/release.sh patch

# Minor release (1.0.0 → 1.1.0)  
./scripts/release.sh minor

# Major release (1.0.0 → 2.0.0)
./scripts/release.sh major
```

Or manually with git:
```bash
git tag v1.0.0
git push origin v1.0.0
```

## 🤝 Want to Help?

1. Fork the repository
2. Make your changes
3. Submit a pull request

We welcome contributions!

## 📄 License

This project is licensed under the MIT License.

## 🙏 Thanks To

- [AssemblyAI](https://www.assemblyai.com/) for the speech recognition
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) for YouTube downloads
- [Cobra](https://github.com/spf13/cobra) for the command-line interface
- [Viper](https://github.com/spf13/viper) for configuration management
