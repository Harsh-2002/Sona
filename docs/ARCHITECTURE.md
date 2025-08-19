# Sona Architecture

## Overview

Sona is a CLI tool built with Go that provides audio-to-text transcription capabilities. It's designed with a modular architecture that separates concerns and makes the codebase maintainable and extensible.

## Architecture Principles

- **Separation of Concerns** - Each package has a single responsibility
- **Dependency Injection** - Dependencies are injected rather than hardcoded
- **Interface-based Design** - Uses interfaces for flexibility and testing
- **Error Handling** - Comprehensive error handling with meaningful messages
- **Configuration Management** - Centralized configuration with multiple sources

## Package Structure

```
pkg/
├── assemblyai/     # External API client
├── config/         # Configuration management
├── transcriber/    # Business logic orchestration
└── youtube/        # YouTube audio download
```

## Package Details

### `pkg/assemblyai`

**Responsibility**: Communication with AssemblyAI REST API

**Key Components**:
- `Client` - HTTP client with authentication
- `TranscribeAudio()` - Main transcription method
- `uploadAudioFile()` - File upload to AssemblyAI
- `submitTranscription()` - Submit transcription request
- `pollTranscription()` - Poll for completion

**Dependencies**: Standard library (`net/http`, `encoding/json`, `mime/multipart`)

### `pkg/config`

**Responsibility**: Configuration management and persistence

**Key Components**:
- `InitConfig()` - Initialize configuration system
- `GetAPIKey()` - Retrieve API key with validation
- `GetOutputPath()` - Get default output directory
- `ConfigCmd` - CLI commands for configuration

**Dependencies**: `github.com/spf13/viper`, `github.com/spf13/cobra`

### `pkg/transcriber`

**Responsibility**: Orchestrating the transcription process

**Key Components**:
- `TranscribeCmd` - Main CLI command
- `processYouTubeVideo()` - Handle YouTube URLs
- `processLocalAudio()` - Handle local files
- `saveTranscript()` - Save results with smart naming

**Dependencies**: All other packages, `github.com/spf13/cobra`

### `pkg/youtube`

**Responsibility**: YouTube audio download using pure Go

**Key Components**:
- `DownloadAudio()` - Download audio from YouTube URL
- `IsYouTubeURL()` - Validate YouTube URLs
- `GetVideoInfo()` - Get video metadata

**Dependencies**: Native yt-dlp binary management

## Data Flow

```
User Input → CLI Command → Transcriber → YouTube/AssemblyAI → Output File
     ↓              ↓           ↓              ↓              ↓
  YouTube URL   transcribe   Determine    Download &     Save to
  or File Path              Source Type   Transcribe     File System
```

## Configuration Flow

```
Environment Variable → Viper → Config File → Runtime
       ↓                ↓         ↓          ↓
  ASSEMBLYAI_API_KEY  Load    ~/.sona/   Use in App
```

## Error Handling Strategy

1. **Input Validation** - Validate user input early
2. **Graceful Degradation** - Provide helpful error messages
3. **Context Preservation** - Include relevant context in errors
4. **User Guidance** - Suggest solutions for common issues

## Testing Strategy

- **Unit Tests** - Test individual package functions
- **Integration Tests** - Test package interactions
- **CLI Tests** - Test command-line interface
- **Mock External APIs** - Use mocks for AssemblyAI and YouTube

## Future Extensibility

### Planned Features
- **Batch Processing** - Process multiple files
- **Audio Preprocessing** - Audio enhancement before transcription
- **Multiple Output Formats** - JSON, SRT, VTT
- **Progress Bars** - Visual download/transcription progress
- **Resume Support** - Resume interrupted downloads

### Architecture Considerations
- **Plugin System** - Support for custom audio sources
- **Queue System** - Background processing of multiple files
- **Caching** - Cache downloaded audio files
- **Metrics** - Performance and usage metrics

## Security Considerations

- **API Key Management** - Secure storage and masking
- **File Validation** - Validate audio file types and sizes
- **URL Validation** - Sanitize YouTube URLs
- **Output Path Validation** - Prevent path traversal attacks

## Performance Considerations

- **Streaming Downloads** - Download audio in chunks
- **Concurrent Processing** - Process multiple files simultaneously
- **Memory Management** - Efficient handling of large audio files
- **Connection Pooling** - Reuse HTTP connections
