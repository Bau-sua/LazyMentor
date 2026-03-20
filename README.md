# LazyMentor

> Your LazyVim learning companion

A TUI installer that sets up an AI mentor for learning LazyVim keybindings and navigation. Works with **OpenCode** and **Claude Code**.

## What is LazyMentor?

LazyMentor is a system prompt that lives inside your AI coding agent. It teaches you LazyVim keybindings through conversation — no code generation, just pure navigation knowledge.

**Example conversation:**
```
You: How do I create a new file?
LazyMentor: Space + e + a (according to your config). If your mapping is different, let me know.

┌─────────────┬──────────────────────────┐
│ Keymap      │ Action                   │
├─────────────┼──────────────────────────┤
│ <leader>e a │ New file (harpoon)       │
│ <leader>nf  │ New file (telescope)     │
└─────────────┴──────────────────────────┘
```

## Features

- 🎯 Teaches LazyVim keybindings through conversation
- 📊 Always uses Markdown tables for keymaps
- 🔒 Never generates code or modifies files
- 🎨 Beautiful TUI installer (with CLI fallback)
- 🔄 Install/uninstall with one command
- 📦 Cross-platform (Linux, macOS, Windows)

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/Bau-sua/LazyMentor/main/install.sh | bash
```

### Manual Install

```bash
# Download the binary for your platform
# Then run:
./lazymint
```

### From Source

```bash
git clone https://github.com/Bau-sua/LazyMentor.git
cd LazyMentor
go build -o lazymint ./cmd/installer/
./lazymint
```

## Usage

### Interactive Mode (TUI)

```bash
./lazymint
```

Navigate with `↑/↓` or `j/k`, press `Enter` to confirm, `q` to quit.

### CLI Mode

```bash
# List detected agents
./lazymint -list

# Install lazymentor
./lazymint -install

# Uninstall lazymentor
./lazymint -uninstall
```

### Silent Mode (for scripts/CI)

```bash
./lazymint  # Automatically uses CLI if no TTY
```

## Supported Agents

| Agent | Config Path |
|-------|-------------|
| **OpenCode** | `~/.config/opencode/` |
| **Claude Code** | `~/.claude/` |

## How It Works

1. **Pre-flight Check** — Confirm you have Neovim open (optional but recommended)
2. **Select Agent** — Choose where to install LazyMentor
3. **Install** — Copies the prompt to your agent's config directory
4. **Learn** — Ask questions about LazyVim keybindings!

## Development

```bash
# Run tests
go test -v ./...

# Build
go build -o lazymint ./cmd/installer/

# Run in development mode (uses local lazymentor.md)
./lazymint
```

## License

MIT
