// Package agents defines the supported AI coding agents
package agents

import (
	"bytes"
	"encoding/json"
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
	// LazyMentorAgentName is the name used in opencode.json
	LazyMentorAgentName = "lazymentor"
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
	ConfigFile  string // Path to the agent's config file (opencode.json, settings.json, etc.)
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
			ConfigFile:  filepath.Join(opencodePath, "opencode.json"),
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
	switch a.Type {
	case OpenCode:
		return a.isOpenCodeInstalled()
	case ClaudeCode:
		return a.isClaudeCodeInstalled()
	default:
		return false
	}
}

func (a Agent) isOpenCodeInstalled() bool {
	// Check if opencode.json has the lazymentor agent entry
	content, err := os.ReadFile(a.ConfigFile)
	if err != nil {
		return false
	}

	// Simple check - look for lazymentor agent in JSON
	return bytes.Contains(content, []byte(`"`+LazyMentorAgentName+`"`))
}

func (a Agent) isClaudeCodeInstalled() bool {
	content, err := os.ReadFile(a.PromptPath)
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

	// 1. Copy lazymentor.md to config directory
	promptPath := filepath.Join(dir, LazyMentorFile)
	if err := os.WriteFile(promptPath, []byte(prompt), 0644); err != nil {
		return fmt.Errorf("failed to write lazymentor.md: %w", err)
	}

	// 2. Modify opencode.json to add the agent
	if err := a.addOpenCodeAgent(); err != nil {
		return fmt.Errorf("failed to add agent to opencode.json: %w", err)
	}

	return nil
}

func (a Agent) addOpenCodeAgent() error {
	// Read existing config
	content, err := os.ReadFile(a.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read opencode.json: %w", err)
	}

	// Create backup
	backupPath := a.ConfigFile + ".lazymentor.backup"
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Parse JSON
	var config map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		return fmt.Errorf("failed to parse opencode.json: %w", err)
	}

	// Ensure agent section exists
	if _, ok := config["agent"]; !ok {
		config["agent"] = map[string]interface{}{}
	}

	agentMap, ok := config["agent"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("agent section is not a valid object")
	}

	// Add lazymentor agent
	agentMap[LazyMentorAgentName] = map[string]interface{}{
		"description": "Your LazyVim learning companion - teaches keybindings through conversation",
		"prompt":      "{file:./" + LazyMentorFile + "}",
		"mode":        "all",
		"tools": map[string]interface{}{
			"read":  true,
			"write": false,
			"edit":  false,
			"bash":  false,
		},
	}

	// Write back
	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Preserve original formatting by adding newline
	output = append(output, '\n')

	if err := os.WriteFile(a.ConfigFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write opencode.json: %w", err)
	}

	return nil
}

func (a Agent) uninstallOpenCode() error {
	// Remove lazymentor.md
	promptPath := filepath.Join(a.ConfigPath, LazyMentorFile)
	if _, err := os.Stat(promptPath); err == nil {
		if err := os.Remove(promptPath); err != nil {
			return fmt.Errorf("failed to remove lazymentor.md: %w", err)
		}
	}

	// Remove agent from opencode.json
	if err := a.removeOpenCodeAgent(); err != nil {
		return fmt.Errorf("failed to remove agent from opencode.json: %w", err)
	}

	return nil
}

func (a Agent) removeOpenCodeAgent() error {
	content, err := os.ReadFile(a.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read opencode.json: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		return fmt.Errorf("failed to parse opencode.json: %w", err)
	}

	// Check if agent section exists
	agentMap, ok := config["agent"].(map[string]interface{})
	if !ok {
		return nil // Nothing to remove
	}

	// Remove lazymentor agent
	delete(agentMap, LazyMentorAgentName)

	// If agent map is empty, remove the section
	if len(agentMap) == 0 {
		delete(config, "agent")
	}

	// Write back
	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	output = append(output, '\n')

	return os.WriteFile(a.ConfigFile, output, 0644)
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
