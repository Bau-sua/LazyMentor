// Package agents defines the supported AI coding agents
package agents

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// AgentType represents a supported AI coding agent
type AgentType string

const (
	OpenCode   AgentType = "opencode"
	ClaudeCode AgentType = "claude-code"
	Other      AgentType = "other"
)

// Agent represents a supported AI coding agent
type Agent struct {
	Name        string
	Type        AgentType
	ConfigPath  string
	Description string
}

// DetectAgent returns a list of detected agents on the system
func DetectAgents() []Agent {
	var agents []Agent

	// OpenCode
	opencodePath := getOpenCodeConfigPath()
	if opencodePath != "" {
		agents = append(agents, Agent{
			Name:        "OpenCode",
			Type:        OpenCode,
			ConfigPath:  opencodePath,
			Description: "OpenCode CLI with opencode.json config",
		})
	}

	// Claude Code
	claudePath := getClaudeConfigPath()
	if claudePath != "" {
		agents = append(agents, Agent{
			Name:        "Claude Code",
			Type:        ClaudeCode,
			ConfigPath:  claudePath,
			Description: "Claude Code CLI with CLAUDE.md",
		})
	}

	return agents
}

func getHomeDir() string {
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

func getOpenCodeConfigPath() string {
	home := getHomeDir()
	var configPath string

	switch runtime.GOOS {
	case "windows":
		configPath = filepath.Join(os.Getenv("APPDATA"), "opencode", "opencode.json")
	default:
		configPath = filepath.Join(home, ".config", "opencode", "opencode.json")
	}

	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	// Fallback to legacy location
	legacyPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath
	}

	return ""
}

func getClaudeConfigPath() string {
	home := getHomeDir()
	configPath := filepath.Join(home, ".claude", "CLAUDE.md")

	if _, err := os.Stat(filepath.Dir(configPath)); err == nil {
		return configPath
	}

	return ""
}

// InstallPrompt installs the lazymentor prompt for the given agent
func (a Agent) InstallPrompt(prompt string) error {
	switch a.Type {
	case OpenCode:
		return a.installOpenCode(prompt)
	case ClaudeCode:
		return a.installClaudeCode(prompt)
	default:
		return fmt.Errorf("unsupported agent type: %s", a.Type)
	}
}

func (a Agent) installOpenCode(prompt string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(a.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// For OpenCode, we copy lazymentor.md to the config directory
	destPath := filepath.Join(dir, "lazymentor.md")
	return os.WriteFile(destPath, []byte(prompt), 0644)
}

func (a Agent) installClaudeCode(prompt string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(a.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if CLAUDE.md already exists
	if _, err := os.Stat(a.ConfigPath); err == nil {
		// File exists, we should ask the user if they want to merge or overwrite
		// For now, we'll append
		existing, _ := os.ReadFile(a.ConfigPath)
		merged := string(existing) + "\n\n---\n\n" + prompt
		return os.WriteFile(a.ConfigPath, []byte(merged), 0644)
	}

	return os.WriteFile(a.ConfigPath, []byte(prompt), 0644)
}
