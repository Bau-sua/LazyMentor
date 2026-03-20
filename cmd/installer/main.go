package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lazymentor/lazymint/internal/agents"
	"github.com/lazymentor/lazymint/internal/embed"
	"github.com/lazymentor/lazymint/internal/nvim"
	"github.com/lazymentor/lazymint/internal/tui"
	"github.com/lazymentor/lazymint/internal/update"
)

const (
	repo = "Bau-sua/LazyMentor"
)

func main() {
	// Parse CLI flags
	installCmd := flag.Bool("install", false, "Install lazymentor to detected agents")
	uninstallCmd := flag.Bool("uninstall", false, "Uninstall lazymentor from detected agents")
	listCmd := flag.Bool("list", false, "List detected agents and installation status")
	versionCmd := flag.Bool("version", false, "Show version information")
	checkCmd := flag.Bool("check-updates", false, "Check for updates")
	updateCmd := flag.Bool("update", false, "Update to the latest version")
	nvimInfoCmd := flag.Bool("nvim-info", false, "Show Neovim configuration info")
	flag.Parse()

	// Version flag (works with any other flags)
	if *versionCmd {
		showVersion()
		return
	}

	// Neovim info
	if *nvimInfoCmd {
		showNvimInfo()
		return
	}

	// Check for updates
	if *checkCmd {
		checkForUpdates()
		return
	}

	// Update
	if *updateCmd {
		performUpdate()
		return
	}

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
	// Check if stdout is a terminal
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	// Check if stdin is a terminal (more reliable for curl | bash)
	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		return true // Assume TTY if we can't check stdin
	}

	return (stdinInfo.Mode() & os.ModeCharDevice) != 0
}

func showVersion() {
	fmt.Printf("LazyMentor %s\n", update.GetCurrentVersion())
	fmt.Printf("Repository: %s\n", repo)
}

func showNvimInfo() {
	cfg := nvim.Detect()
	fmt.Println(cfg.FormatForPrompt())
}

func checkForUpdates() {
	fmt.Println("Checking for updates...")
	fmt.Println()

	result, err := update.CheckForUpdates(repo)
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current version: %s\n", result.CurrentVersion)
	fmt.Printf("Latest version:  %s\n", result.LatestVersion)
	fmt.Println()

	if result.UpdateNeeded {
		fmt.Println("A new version is available!")
		fmt.Printf("Download: %s\n", result.DownloadURL)
		fmt.Println()
		fmt.Println("To update, run: lazymint --update")
	} else {
		fmt.Println("You're running the latest version.")
	}
}

func performUpdate() {
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not determine binary path: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Checking for updates...")
	fmt.Println()

	result, err := update.CheckForUpdates(repo)
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	if !result.UpdateNeeded {
		fmt.Println("You're already running the latest version!")
		return
	}

	fmt.Printf("Updating from %s to %s...\n", result.CurrentVersion, result.LatestVersion)
	fmt.Println()

	if err := update.DownloadAndReplace(result.DownloadURL, binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Update complete!")
	fmt.Println("Run 'lazymint --version' to verify.")
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
	fmt.Println("LazyMentor CLI")
	fmt.Printf("Version: %s\n", update.GetCurrentVersion())
	fmt.Println()

	// Show Neovim info
	cfg := nvim.Detect()
	if cfg.Installed {
		fmt.Println("Neovim:")
		if cfg.IsLazyVim {
			fmt.Printf("  ✓ LazyVim %s detected\n", cfg.Version)
		} else {
			fmt.Printf("  ✓ Neovim %s detected\n", cfg.Version)
		}
		fmt.Printf("  Config: %s\n", cfg.ConfigDir)
		if len(cfg.Plugins) > 0 {
			fmt.Printf("  Plugins: %s\n", joinStrings(cfg.Plugins, ", "))
		}
		fmt.Println()
	}

	// Show agents
	fmt.Println("Agents:")
	if len(detectedAgents) == 0 {
		fmt.Println("  (none detected)")
	}
	for _, agent := range detectedAgents {
		installed := "not installed"
		if agent.IsInstalled() {
			installed = "installed"
		}
		fmt.Printf("  • %s - %s\n", agent.Name, installed)
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

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
