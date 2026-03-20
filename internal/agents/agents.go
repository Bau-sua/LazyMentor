// Package agents defines the supported AI coding agents
package agents

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// LazyMentorFile is the name of the lazymentor prompt file
	LazyMentorFile = "lazymentor.md"
	// LazyMentorMarker is used to identify if lazymentor is installed
	LazyMentorMarker = "# LazyMentor - System Prompt"
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
	PromptPath  string // Path where lazymentor.md will be installed
	Description string
}

// DetectAgents returns a list of detected agents on the system
func DetectAgents() []Agent {
	var agents []Agent

	// OpenCode
	opencodePath := getOpenCodeConfigDir()
	if opencodePath != "" {
		agents = append(agents, Agent{
			Name:        "OpenCode",
			Type:        OpenCode,
			ConfigPath:  opencodePath,
			PromptPath:  filepath.Join(opencodePath, LazyMentorFile),
			Description: "OpenCode CLI with opencode.json config",
		})
	}

	// Claude Code
	claudePath := getClaudeConfigDir()
	if claudePath != "" {
		agents = append(agents, Agent{
			Name:        "Claude Code",
			Type:        ClaudeCode,
			ConfigPath:  claudePath,
			PromptPath:  filepath.Join(claudePath, "CLAUDE.md"),
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

func getOpenCodeConfigDir() string {
	home := getHomeDir()
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "opencode")
	default:
		configDir = filepath.Join(home, ".config", "opencode")
	}

	// Check if directory exists and has opencode.json
	configFile := filepath.Join(configDir, "opencode.json")
	if _, err := os.Stat(configFile); err == nil {
		return configDir
	}

	return ""
}

func getClaudeConfigDir() string {
	home := getHomeDir()
	configDir := filepath.Join(home, ".claude")

	// Check if directory exists
	if _, err := os.Stat(configDir); err == nil {
		return configDir
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

// UninstallPrompt removes the lazymentor prompt from the agent
func (a Agent) UninstallPrompt() error {
	switch a.Type {
	case OpenCode:
		return a.uninstallOpenCode()
	case ClaudeCode:
		return a.uninstallClaudeCode()
	default:
		return fmt.Errorf("unsupported agent type: %s", a.Type)
	}
}

// IsInstalled checks if lazymentor is installed for this agent
func (a Agent) IsInstalled() bool {
	var promptPath string

	switch a.Type {
	case OpenCode:
		promptPath = filepath.Join(a.ConfigPath, LazyMentorFile)
	case ClaudeCode:
		promptPath = a.PromptPath
	default:
		return false
	}

	content, err := os.ReadFile(promptPath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), LazyMentorMarker)
}

func (a Agent) installOpenCode(prompt string) error {
	dir := a.ConfigPath
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	destPath := filepath.Join(dir, LazyMentorFile)
	return os.WriteFile(destPath, []byte(prompt), 0644)
}

func (a Agent) uninstallOpenCode() error {
	promptPath := filepath.Join(a.ConfigPath, LazyMentorFile)

	if _, err := os.Stat(promptPath); os.IsNotExist(err) {
		return fmt.Errorf("lazymentor is not installed for OpenCode")
	}

	return os.Remove(promptPath)
}

func (a Agent) installClaudeCode(prompt string) error {
	dir := a.ConfigPath
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	promptPath := filepath.Join(dir, "CLAUDE.md")

	// Check if CLAUDE.md already exists
	if _, err := os.Stat(promptPath); err == nil {
		// File exists - extract and replace only the lazymentor section
		existing, _ := os.ReadFile(promptPath)
		cleaned := removeLazyMentorSection(string(existing))
		merged := strings.TrimSpace(cleaned) + "\n\n" + prompt
		return os.WriteFile(promptPath, []byte(merged), 0644)
	}

	// No existing CLAUDE.md, create with lazymentor
	return os.WriteFile(promptPath, []byte(prompt), 0644)
}

func (a Agent) uninstallClaudeCode() error {
	promptPath := a.PromptPath

	if _, err := os.Stat(promptPath); os.IsNotExist(err) {
		return fmt.Errorf("lazymentor is not installed for Claude Code")
	}

	existing, err := os.ReadFile(promptPath)
	if err != nil {
		return err
	}

	cleaned := removeLazyMentorSection(string(existing))
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		// Nothing left, remove the file
		return os.Remove(promptPath)
	}

	// Write back the cleaned content
	return os.WriteFile(promptPath, []byte(cleaned), 0644)
}

// removeLazyMentorSection removes the lazymentor section from Claude.md content
func removeLazyMentorSection(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inSection := false

	for _, line := range lines {
		if strings.Contains(line, LazyMentorMarker) {
			inSection = true
			continue
		}

		// End of lazymentor section (next heading or ---)
		if inSection && (strings.HasPrefix(line, "#") || strings.HasPrefix(line, "---")) {
			inSection = false
		}

		if !inSection {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// ContainsLazyMentor checks if content contains the lazymentor prompt
func ContainsLazyMentor(content []byte) bool {
	return bytes.Contains(content, []byte(LazyMentorMarker))
}
