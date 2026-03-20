package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lazymentor/lazymint/internal/agents"
	"github.com/lazymentor/lazymint/internal/embed"
	"github.com/lazymentor/lazymint/internal/tui"
)

func main() {
	// Parse CLI flags
	installCmd := flag.Bool("install", false, "Install lazymentor to detected agents")
	uninstallCmd := flag.Bool("uninstall", false, "Uninstall lazymentor from detected agents")
	listCmd := flag.Bool("list", false, "List detected agents and installation status")
	flag.Parse()

	// CLI mode
	if *installCmd || *uninstallCmd || *listCmd {
		runCLI(*installCmd, *uninstallCmd, *listCmd)
		return
	}

	// TUI mode (if terminal supports it)
	if isTTY() {
		if err := tui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Fallback: silent install
	runSilentInstall()
}

func isTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func runCLI(install, uninstall, list bool) {
	detectedAgents := agents.DetectAgents()

	if list {
		listAgents(detectedAgents)
		return
	}

	if len(detectedAgents) == 0 {
		fmt.Println("No supported agents detected.")
		fmt.Println("Supported agents: OpenCode, Claude Code")
		if list {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if uninstall {
		uninstallAll(detectedAgents)
		return
	}

	if install {
		installAll(detectedAgents)
		return
	}
}

func listAgents(detectedAgents []agents.Agent) {
	if len(detectedAgents) == 0 {
		fmt.Println("No supported agents detected.")
		return
	}

	fmt.Println("Detected agents:")
	fmt.Println()
	for _, agent := range detectedAgents {
		installed := "not installed"
		if agent.IsInstalled() {
			installed = "installed"
		}
		fmt.Printf("  • %s (%s) - %s\n", agent.Name, agent.ConfigPath, installed)
	}
}

func installAll(detectedAgents []agents.Agent) {
	prompt := loadLocalPrompt()
	if prompt == "" {
		prompt = embed.LazyMentorPrompt
	}

	for _, agent := range detectedAgents {
		if agent.IsInstalled() {
			fmt.Printf("Skipping %s (already installed)\n", agent.Name)
			continue
		}

		fmt.Printf("Installing to %s...\n", agent.Name)
		if err := agent.InstallPrompt(prompt); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing to %s: %v\n", agent.Name, err)
			continue
		}
		fmt.Printf("✓ Installed to %s\n", agent.Name)
	}
}

func uninstallAll(detectedAgents []agents.Agent) {
	for _, agent := range detectedAgents {
		if !agent.IsInstalled() {
			fmt.Printf("Skipping %s (not installed)\n", agent.Name)
			continue
		}

		fmt.Printf("Uninstalling from %s...\n", agent.Name)
		if err := agent.UninstallPrompt(); err != nil {
			fmt.Fprintf(os.Stderr, "Error uninstalling from %s: %v\n", agent.Name, err)
			continue
		}
		fmt.Printf("✓ Uninstalled from %s\n", agent.Name)
	}
}

func runSilentInstall() {
	fmt.Println("LazyMentor Installer (silent mode)")
	fmt.Println()

	prompt := loadLocalPrompt()
	if prompt == "" {
		prompt = embed.LazyMentorPrompt
		fmt.Println("Using embedded prompt")
	} else {
		fmt.Println("Using local prompt: lazymentor.md")
	}

	detectedAgents := agents.DetectAgents()

	if len(detectedAgents) == 0 {
		fmt.Println("No supported agents detected.")
		fmt.Println("Supported agents: OpenCode, Claude Code")
		os.Exit(1)
	}

	fmt.Printf("Detected %d agent(s):\n", len(detectedAgents))
	for _, agent := range detectedAgents {
		installed := ""
		if agent.IsInstalled() {
			installed = " (already installed)"
		}
		fmt.Printf("  - %s (%s)%s\n", agent.Name, agent.ConfigPath, installed)
	}

	// Install to first detected agent
	agent := detectedAgents[0]
	if agent.IsInstalled() {
		fmt.Printf("\nSkipping %s (already installed)\n", agent.Name)
		os.Exit(0)
	}

	fmt.Printf("\nInstalling to %s...\n", agent.Name)
	if err := agent.InstallPrompt(prompt); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Installation complete!")
	fmt.Printf("Prompt installed to: %s\n", agent.PromptPath)
	fmt.Println("\nRestart your agent to start learning LazyVim!")
}

func loadLocalPrompt() string {
	data, err := os.ReadFile("lazymentor.md")
	if err != nil {
		return ""
	}
	return string(data)
}
