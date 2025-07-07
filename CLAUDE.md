# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go application called "basay" that serves as a text-to-speech "roast" tool. It takes text input (typically Japanese) and converts it into harsh/sarcastic responses, with optional OpenAI integration for both text generation and audio output.

## Architecture

The application has three main components:

1. **Local Roast Generation** (`localRoast` function): Pattern-matching based responses for common Japanese phrases
2. **ChatGPT Integration** (`chatRoast` function): Uses OpenAI's GPT-4o-mini model to generate creative roasts
3. **Text-to-Speech** (`speakTTS` function): Converts generated text to audio using OpenAI's TTS API and plays it via ffplay

The application gracefully degrades from OpenAI integration → local responses if API key is missing or API calls fail.

## Development Commands

### Build
```bash
go build -o roast main.go
```

### Run
```bash
# With OpenAI integration
OPENAI_API_KEY=sk-xxx echo "遅れてすみません" | ./roast

# Local-only mode (no API key)
echo "遅れてすみません" | ./roast
```

### Dependencies
```bash
go mod tidy
```

## Dependencies

- `github.com/openai/openai-go` - OpenAI API client
- Requires `ffplay` (from FFmpeg) for audio playback

## Key Files

- `main.go` - Single-file application containing all functionality
- `go.mod` - Go module definition with OpenAI dependency
- `README.md` - Basic project description

## Environment Variables

- `OPENAI_API_KEY` - Required for ChatGPT and TTS functionality (optional for local-only mode)

## Usage Pattern

The application reads from stdin and outputs either:
1. Audio playback (if OpenAI integration works)
2. Text output to stdout (fallback)

Input is typically Japanese text expressing common situations (lateness, excuses, etc.) that trigger specific roast responses.