# Sona - Audio Transcription Tool (CLI)

Turn audio files and YouTube videos into text with just a few commands. Sona uses AssemblyAI's speech recognition to give you accurate transcripts quickly and easily.

## 🚀 Quick Install

### Homebrew Installation (Recommended for macOS/Linux)

**Install via Homebrew:**
```bash
# Add the Sona tap
brew tap harsh-2002/sona

# Install Sona
brew install sona
```



### One-liner Installation (Alternative)

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/Harsh-2002/Sona/main/install.sh | sudo sh
```

**With wget:**
```bash
wget -qO- https://raw.githubusercontent.com/Harsh-2002/Sona/main/install.sh | sudo sh
```

### Manual Installation

1. **Download the installer:**
   ```bash
   curl -O https://raw.githubusercontent.com/Harsh-2002/Sona/main/install.sh
   chmod +x install.sh
   ```

2. **Run the installer (requires sudo):**
   ```bash
   sudo ./install.sh
   ```

### Uninstalling Sona

**Remove Sona completely:**
```bash
sudo ./install.sh --uninstall
```

The uninstaller will:
- Remove the Sona binary
- Ask if you want to remove auto-installed dependencies (yt-dlp, FFmpeg)
- Clean up configuration files

**Note**: The installer requires root privileges to install Sona system-wide in `/usr/local/bin/`, making it accessible from anywhere on your system.

## 📥 Direct Downloads

If you prefer to download manually, here are the direct links for each platform:

### Linux
- **AMD64**: [sona-linux-amd64](https://s3.srvr.site/artifact/sona/sona-linux-amd64)
- **ARM64**: [sona-linux-arm64](https://s3.srvr.site/artifact/sona/sona-linux-arm64)

### macOS
- **Intel**: [sona-darwin-amd64](https://s3.srvr.site/artifact/sona/sona-darwin-amd64)
- **Apple Silicon**: [sona-darwin-arm64](https://s3.srvr.site/artifact/sona/sona-darwin-arm64)

### Windows
- **AMD64**: [sona-windows-amd64.exe](https://s3.srvr.site/artifact/sona/sona-windows-amd64.exe)
- **ARM64**: [sona-windows-arm64.exe](https://s3.srvr.site/artifact/sona/sona-windows-arm64.exe)

## ✨ What Sona Does

- **Audio Files** - Convert any audio file (MP3, WAV, M4A, etc.) to text
- **YouTube Videos** - Download and transcribe YouTube videos automatically
- **Smart AI** - Uses the latest speech recognition models for best accuracy
- **Easy Setup** - Simple configuration with your API key
- **Flexible Output** - Save transcripts wherever you want with smart naming
- **Secure** - Your API keys are encrypted and safe
- **Interactive** - Step-by-step guidance when you need it

## 🔧 What You Need

- **AssemblyAI API Key** - Get one at [assemblyai.com](https://assemblyai.com)
- **Internet Connection** - For API calls and YouTube downloads
- **Storage Space** - For temporary audio files and transcripts

The installer automatically handles dependencies like `yt-dlp` and `FFmpeg` for you. After installation, run `sona install` to set up the required dependencies for your platform.

## 🏗️ How Sona Works

Sona is built with Go and uses these main parts:

- **CLI Interface** - Easy commands for transcription
- **YouTube Downloader** - Gets audio from YouTube videos
- **Audio Converter** - Changes audio to the right format
- **AssemblyAI Client** - Sends audio to AI for transcription
- **Configuration Manager** - Keeps your settings safe

## 🚀 How Sona is Built

- **Go 1.24+** - Fast, reliable programming language
- **Cobra** - Professional command-line interface
- **Viper** - Smart configuration management
- **AssemblyAI API** - Industry-leading speech recognition

## 🔨 Building from Source

If you want to build Sona yourself:

```bash
git clone https://github.com/Harsh-2002/Sona.git
cd Sona
go build -o sona cmd/sona/main.go
```

## 📖 How to Use Sona

### Basic Commands

**Transcribe a local audio file:**
```bash
sona transcribe audio.mp3
```

**Transcribe a YouTube video:**
```bash
sona transcribe https://youtube.com/watch?v=VIDEO_ID
```

**Start interactive mode:**
```bash
sona interactive
```

**Install dependencies:**
```bash
sona install
```

**Manage your settings:**
```bash
sona config
```

### Command Options

- `--output` - Save transcript to specific file
- `--model` - Choose AI model (default: best)
- `--language` - Set audio language (auto-detected by default)

## 🤖 AI Models

Sona uses AssemblyAI's latest models:

- **Best** - Highest accuracy (default)
- **Nano** - Fastest processing
- **Base** - Balanced speed and accuracy

## 📁 Where Transcripts Go

By default, transcripts are saved to:
- **Current directory** - `./transcript.txt`
- **Custom path** - Use `--output` flag
- **Smart naming** - Based on original filename

## ⚙️ Settings

Sona stores your settings in `~/.sona/config.toml`:

- **API Key** - Your AssemblyAI access key
- **Default Model** - Preferred AI model
- **Output Directory** - Where to save transcripts
- **Language** - Default audio language

## 🔒 Keeping Your Data Safe

- **API Keys** - Encrypted with AES-256-GCM
- **Local Storage** - Files stay on your device
- **No Data Collection** - Sona doesn't track your usage
- **Secure API** - HTTPS for all communications

## 🚨 When Things Go Wrong

### Common Problems

**"API key not found"**
- Run `sona config` to set your key
- Check environment variable `ASSEMBLYAI_API_KEY`

**"yt-dlp not found"**
- Run `sona install` to install dependencies
- Manual install: Dependencies must be installed using `sona install`

**"FFmpeg not found"**
- Run `sona install` to install dependencies
- Manual install: Dependencies must be installed using `sona install`

**"xz utility not found" or "failed to extract FFmpeg archive"**
- Install xz-utils: `sudo apt-get install xz-utils` (Ubuntu/Debian) or `sudo yum install xz` (CentOS/RHEL)
- This is required for extracting FFmpeg archives on some Linux distributions

**"Permission denied"**
- Run installer with `sudo ./install.sh`
- For uninstall: `sudo ./install.sh --uninstall`
- Check file permissions

### Debugging

**Check the log file:**
```bash
cat ~/.sona/sona.log
```

The log file contains detailed information about:
- Dependency installation
- YouTube downloads  
- Audio processing
- Transcription steps
- Error details

**Location:** `~/.sona/sona.log`

**macOS Note:** On macOS, Sona automatically installs both `ffmpeg` and `ffprobe` from evermeet.cx, which are required for YouTube audio extraction.

**Path Consistency:** On all Unix-like systems (Linux, macOS, BSD), dependencies are installed to `~/bin/` for consistent behavior across platforms.

**System Requirements:** Some Linux distributions may require the `xz-utils` package for FFmpeg installation. If you encounter extraction errors, install it with:
- **Ubuntu/Debian**: `sudo apt-get install xz-utils`
- **CentOS/RHEL**: `sudo yum install xz`
- **Alpine**: `apk add xz`

### Updating Sona

**If installed via Homebrew:**
```bash
brew upgrade sona
```

**If installed via installer:**
```bash
sudo ./install.sh
```

The installer automatically detects if you have an older version and updates it. No need for separate upgrade commands!

### Getting Help

- **Check logs** - Look for error messages
- **Verify API key** - Test with `sona config`
- **Check internet** - Ensure connectivity
- **Update Sona** - Get the latest version

## 🔄 Automated Builds

Sona automatically builds and uploads new binaries on every push to `main` or `master` branches. The GitHub Actions workflow:

- **Builds for all platforms** - Linux, macOS, Windows (AMD64/ARM64)
- **Uploads to MinIO** - Public bucket at `https://s3.srvr.site/artifact/sona/`
- **No manual releases** - Everything happens automatically
- **Always up-to-date** - Latest code = latest binaries

**Note**: The `sona` folder in the MinIO bucket is created automatically if it doesn't exist.

## 🚀 What's Coming Next

- **More AI Models** - Additional speech recognition options like OpenAI Whispher
- **Batch Processing** - Handle multiple files at once
- **Export Formats** - JSON, SRT, VTT support
- **Real-time Transcription** - Live audio processing
- **Slack notifications** - Get notified when transcription completes
- **Notion/Notion-like** - Export to note-taking apps

## 🤝 Want to Help?

Contributions are welcome! Here's how:

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Test thoroughly**
5. **Submit a pull request**

## 🙏 Thanks To

- **AssemblyAI** - For amazing speech recognition
- **Go Community** - For excellent tools and libraries
- **Open Source** - For making this possible

## ❓ Need Help?

- **GitHub Issues** - Report bugs and request features
- **Documentation** - Check this README first
- **Community** - Join discussions and get help

---

**Made with ❤️ for easy audio transcription**
