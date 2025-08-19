# Basic Usage Examples

## Quick Start

### 1. Set API Key
```bash
# Set your AssemblyAI API key
export ASSEMBLYAI_API_KEY="your_api_key_here"
```

### 2. Transcribe YouTube Video
```bash
./sona transcribe "https://youtube.com/watch?v=dQw4w9WgXcQ"
```

### 3. Transcribe Local Audio File
```bash
./sona transcribe "./meeting_recording.mp3"
```

## Advanced Usage

### Custom Output Path
```bash
./sona transcribe "video.mp4" --output ./my_transcript.txt
```

### Different Speech Model
```bash
# Use best quality model
./sona transcribe "audio.mp3" --model best

# Use fastest model
./sona transcribe "audio.mp3" --model nano
```

### Configuration Management
```bash
# Show current config
./sona config show

# Set API key via config
./sona config set api_key "your_new_key_here"
```

## Output Examples

### YouTube Video Output
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

### Local File Output
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
