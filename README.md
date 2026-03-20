# LazyMentor

> Your LazyVim learning companion

A TUI installer that sets up an AI mentor for learning LazyVim keybindings and navigation. Works with **OpenCode** and **Claude Code**.

---

## What is LazyMentor?

LazyMentor is a system prompt that lives inside your AI coding agent. It teaches you LazyVim keybindings through conversation — **no code generation, no file modifications, just pure navigation knowledge**.

Instead of searching documentation or watching tutorials, you can now ask questions naturally:

```
You: How do I create a new file?
LazyMentor: Space + e + a (according to your config). 
            If your mapping is different, let me know.

┌─────────────┬──────────────────────────┐
│ Keymap      │ Action                   │
├─────────────┼──────────────────────────┤
│ <leader>e a │ New file (harpoon)       │
│ <leader>nf  │ New file (telescope)     │
└─────────────┴──────────────────────────┘
```

### Key principles

- **Never generates code** — only explains keybindings
- **Never modifies files** — only teaches you how to do it yourself
- **Always uses tables** — clear, visual keymap references
- **Respects your config** — if you paste your `init.lua`, it uses that as truth

---

## Features

- 🎯 Teaches LazyVim keybindings through natural conversation
- 📊 Always responds with Markdown tables for keymaps
- 🔒 **Zero risk**: never generates code or touches files
- 🎨 Beautiful TUI installer (with CLI fallback)
- 🔄 Install/uninstall in one command
- 📦 Cross-platform: Linux, macOS, Windows

---

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/Bau-sua/LazyMentor/main/install.sh | bash
```

This will:
1. Download the correct binary for your OS/architecture
2. Install to `~/.local/bin/` or `/usr/local/bin/`
3. Launch the installer automatically

### Manual Install

**Step 1: Download a binary**

Download from [Releases](https://github.com/Bau-sua/LazyMentor/releases) for your platform:

| Platform | Binary |
|----------|--------|
| Linux amd64 | `lazymint-linux-amd64` |
| Linux arm64 | `lazymint-linux-arm64` |
| macOS amd64 | `lazymint-darwin-amd64` |
| macOS arm64 (Apple Silicon) | `lazymint-darwin-arm64` |
| Windows | `lazymint-windows-amd64.exe` |

**Step 2: Make it executable**

```bash
chmod +x lazymint-*
```

**Step 3: Run the installer**

```bash
./lazymint
```

### Install from Source

```bash
git clone https://github.com/Bau-sua/LazyMentor.git
cd LazyMentor
go build -o lazymint ./cmd/installer/
./lazymint
```

---

## Usage

### Interactive Mode (TUI)

Launch the installer with:

```bash
./lazymint
```

Navigate with:
- `↑/↓` or `j/k` — move between options
- `Enter` — confirm selection
- `q` — quit

### CLI Mode

For scripts, CI, or quick operations:

```bash
# List detected agents and Neovim info
./lazymint --list

# Install lazymentor to all detected agents
./lazymint --install

# Uninstall lazymentor from all agents
./lazymint --uninstall

# Check for updates
./lazymint --check-updates

# Update to latest version
./lazymint --update

# Show Neovim configuration info
./lazymint --nvim-info

# Show version
./lazymint --version
```

### Silent Mode

When run without a terminal (e.g., `curl | bash`), it automatically uses CLI mode and installs to the first detected agent.

---

## How It Works

### For OpenCode

1. Copies `lazymentor.md` to `~/.config/opencode/`
2. Adds a `lazymentor` agent entry to `opencode.json`
3. The agent appears when you press **Tab** in OpenCode

### For Claude Code

1. Adds the prompt to `~/.claude/CLAUDE.md`
2. Claude Code loads it automatically on startup
3. LazyMentor becomes part of your global instructions

### Uninstall

The uninstall removes:
- The `lazymentor.md` file
- The agent entry from `opencode.json` (OpenCode)
- The lazymentor section from `CLAUDE.md` (Claude Code)
- Your existing configurations are preserved

---

## Supported Agents

| Agent | Config Location | How It Installs |
|-------|----------------|-----------------|
| **OpenCode** | `~/.config/opencode/` | Adds `lazymentor` agent to `opencode.json` |
| **Claude Code** | `~/.claude/` | Appends to `~/.claude/CLAUDE.md` |

---

## Example Questions

Once installed, try asking:

- "How do I create a new file?"
- "Show me how to switch between buffers"
- "How do I open telescope?"
- "Explain splits and how to navigate windows"
- "How do I use harpoon to jump between files?"
- "What is the visual mode and how do I select text?"

---

## Development

```bash
# Run tests
go test -v ./...

# Build for current platform
go build -o lazymint ./cmd/installer/

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o lazymint-linux-amd64 ./cmd/installer/
GOOS=darwin GOARCH=arm64 go build -o lazymint-darwin-arm64 ./cmd/installer/
```

---

## Changelog

### [v0.2.0] - 2026-03-20

#### Added
- **Auto-update**: `--check-updates` and `--update` commands
- **Neovim detection**: Detect version, LazyVim, leader key, and plugins
- **`--version`**: Show version information
- **`--nvim-info`**: Display Neovim configuration details
- **`--list`**: Enhanced with Neovim info display

#### Technical
- Added `internal/nvim` package for Neovim detection
- Added `internal/update` package for GitHub releases
- 26 passing unit tests

### [v0.1.0] - 2026-03-20

#### Added
- TUI installer with bubbletea
- CLI fallback for scripts and CI
- Support for OpenCode (agent registration in `opencode.json`)
- Support for Claude Code (appends to `CLAUDE.md`)
- Install/uninstall/list commands
- Pre-flight checklist (optional)
- 21 passing unit tests
- Cross-platform binaries (Linux, macOS, Windows)

---

## License

MIT
