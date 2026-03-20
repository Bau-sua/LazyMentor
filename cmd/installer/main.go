package main

import (
	"fmt"
	"os"

	"github.com/lazymentor/lazymint/internal/agents"
	"github.com/lazymentor/lazymint/internal/embed"
	"github.com/lazymentor/lazymint/internal/tui"
)

// Silent install mode for non-TTY environments (CI, scripts, etc.)
const silentInstallEnv = "LAZYMINT_SILENT"

func main() {
	// Check for silent install mode
	if os.Getenv(silentInstallEnv) == "1" || !isTTY() {
		runSilentInstall()
		return
	}

	// Normal TUI mode
	if err := tui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func isTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func runSilentInstall() {
	fmt.Println("LazyMentor Installer (silent mode)")
	fmt.Println()

	// Load prompt (local first, then embedded)
	prompt := loadLocalPrompt()
	if prompt == "" {
		prompt = embed.LazyMentorPrompt
		fmt.Println("Using embedded prompt")
	} else {
		fmt.Println("Using local prompt: lazymentor.md")
	}

	// Detect agents
	detectedAgents := agents.DetectAgents()

	if len(detectedAgents) == 0 {
		fmt.Println("No supported agents detected.")
		fmt.Println("Supported agents: OpenCode, Claude Code")
		os.Exit(1)
	}

	fmt.Printf("Detected %d agent(s):\n", len(detectedAgents))
	for _, agent := range detectedAgents {
		fmt.Printf("  - %s (%s)\n", agent.Name, agent.ConfigPath)
	}

	// Install to first detected agent
	agent := detectedAgents[0]
	fmt.Printf("\nInstalling to %s...\n", agent.Name)

	if err := agent.InstallPrompt(prompt); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Installation complete!")
	fmt.Printf("Prompt installed to: %s\n", agent.ConfigPath)
	fmt.Println("\nRestart your agent to start learning LazyVim!")
}

func loadLocalPrompt() string {
	data, err := os.ReadFile("lazymentor.md")
	if err != nil {
		return ""
	}
	return string(data)
}
