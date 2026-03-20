package agents

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectAgents(t *testing.T) {
	agents := DetectAgents()

	if len(agents) == 0 {
		t.Log("No agents detected (expected in test environment)")
	}

	for _, agent := range agents {
		if agent.Name == "" {
			t.Error("Agent name should not be empty")
		}
		if agent.ConfigPath == "" {
			t.Error("Agent config path should not be empty")
		}
		if agent.Type != OpenCode && agent.Type != ClaudeCode {
			t.Errorf("Unexpected agent type: %s", agent.Type)
		}
	}
}

func TestAgentIsInstalled(t *testing.T) {
	agents := DetectAgents()

	for _, agent := range agents {
		// Should not be installed initially (we uninstalled)
		installed := agent.IsInstalled()
		t.Logf("Agent %s installed status: %v", agent.Name, installed)
	}
}

func TestRemoveLazyMentorSection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty content",
			input:    "",
			expected: "",
		},
		{
			name:     "Only lazymentor content",
			input:    "# LazyMentor - System Prompt\nSome content here",
			expected: "",
		},
		{
			name:     "Content before lazymentor",
			input:    "# My Config\nSome config\n\n# LazyMentor - System Prompt\nMentor content",
			expected: "# My Config\nSome config\n",
		},
		{
			name:     "Content after lazymentor",
			input:    "# LazyMentor - System Prompt\nMentor content\n\n# Other Section\nMore content",
			expected: "# Other Section\nMore content",
		},
		{
			name:     "No lazymentor section",
			input:    "# My Config\nSome config\n\n# Other\nMore",
			expected: "# My Config\nSome config\n\n# Other\nMore",
		},
		{
			name:     "With separator",
			input:    "Some content\n\n---\n\n# LazyMentor - System Prompt\nMentor\n\n---\n\nMore content",
			expected: "Some content\n\n---\n\n---\n\nMore content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeLazyMentorSection(tt.input)
			if result != tt.expected {
				t.Errorf("removeLazyMentorSection() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestContainsLazyMentor(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "Contains marker",
			content:  []byte("# LazyMentor - System Prompt\nSome content"),
			expected: true,
		},
		{
			name:     "Does not contain marker",
			content:  []byte("# Other Prompt\nSome content"),
			expected: false,
		},
		{
			name:     "Empty content",
			content:  []byte(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsLazyMentor(tt.content)
			if result != tt.expected {
				t.Errorf("ContainsLazyMentor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInstallPromptToOpenCode(t *testing.T) {
	// Create a temp directory for testing
	tmpDir := t.TempDir()

	agent := Agent{
		Name:        "OpenCode Test",
		Type:        OpenCode,
		ConfigPath:  tmpDir,
		PromptPath:  filepath.Join(tmpDir, LazyMentorFile),
		ConfigFile:  filepath.Join(tmpDir, "opencode.json"),
		Description: "Test agent",
	}

	// Create a minimal opencode.json
	os.WriteFile(agent.ConfigFile, []byte(`{"test": true}`), 0644)

	testPrompt := "# LazyMentor - System Prompt\nTest content"

	err := agent.InstallPrompt(testPrompt)
	if err != nil {
		t.Fatalf("InstallPrompt() error = %v", err)
	}

	// Verify lazymentor.md was created
	installedPath := filepath.Join(tmpDir, LazyMentorFile)
	content, err := os.ReadFile(installedPath)
	if err != nil {
		t.Fatalf("Failed to read installed file: %v", err)
	}

	if string(content) != testPrompt {
		t.Errorf("Installed content = %q, want %q", string(content), testPrompt)
	}

	// Verify IsInstalled returns true
	if !agent.IsInstalled() {
		t.Error("IsInstalled() should return true after installation")
	}

	// Verify opencode.json was modified with lazymentor agent
	jsonContent, err := os.ReadFile(agent.ConfigFile)
	if err != nil {
		t.Fatalf("Failed to read opencode.json: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(jsonContent, &config); err != nil {
		t.Fatalf("Failed to parse opencode.json: %v", err)
	}

	agentSection, ok := config["agent"].(map[string]interface{})
	if !ok {
		t.Fatal("Agent section not found in opencode.json")
	}

	if _, ok := agentSection[LazyMentorAgentName]; !ok {
		t.Error("lazymentor agent not found in opencode.json")
	}
}

func TestOpenCodeAgentJSONModification(t *testing.T) {
	tmpDir := t.TempDir()

	agent := Agent{
		Name:        "OpenCode Test",
		Type:        OpenCode,
		ConfigPath:  tmpDir,
		PromptPath:  filepath.Join(tmpDir, LazyMentorFile),
		ConfigFile:  filepath.Join(tmpDir, "opencode.json"),
		Description: "Test agent",
	}

	// Create opencode.json with existing agents
	existingConfig := `{
  "agent": {
    "build": {
      "mode": "primary",
      "prompt": "Build agent"
    }
  }
}`
	os.WriteFile(agent.ConfigFile, []byte(existingConfig), 0644)

	testPrompt := "# LazyMentor - System Prompt\nTest content"
	if err := agent.InstallPrompt(testPrompt); err != nil {
		t.Fatalf("InstallPrompt() error = %v", err)
	}

	// Verify existing agent is preserved
	jsonContent, _ := os.ReadFile(agent.ConfigFile)
	var config map[string]interface{}
	json.Unmarshal(jsonContent, &config)

	agentSection := config["agent"].(map[string]interface{})

	if _, ok := agentSection["build"]; !ok {
		t.Error("Existing 'build' agent was removed")
	}

	if _, ok := agentSection[LazyMentorAgentName]; !ok {
		t.Error("lazymentor agent not added")
	}
}

func TestUninstallPromptFromOpenCode(t *testing.T) {
	tmpDir := t.TempDir()

	agent := Agent{
		Name:        "OpenCode Test",
		Type:        OpenCode,
		ConfigPath:  tmpDir,
		PromptPath:  filepath.Join(tmpDir, LazyMentorFile),
		ConfigFile:  filepath.Join(tmpDir, "opencode.json"),
		Description: "Test agent",
	}

	// Create a minimal opencode.json
	os.WriteFile(agent.ConfigFile, []byte(`{"test": true}`), 0644)

	// Install first
	testPrompt := "# LazyMentor - System Prompt\nTest content"
	if err := agent.InstallPrompt(testPrompt); err != nil {
		t.Fatalf("InstallPrompt() error = %v", err)
	}

	// Uninstall
	err := agent.UninstallPrompt()
	if err != nil {
		t.Fatalf("UninstallPrompt() error = %v", err)
	}

	// Verify file was removed
	installedPath := filepath.Join(tmpDir, LazyMentorFile)
	if _, err := os.Stat(installedPath); !os.IsNotExist(err) {
		t.Error("File should not exist after uninstall")
	}

	// Verify IsInstalled returns false
	if agent.IsInstalled() {
		t.Error("IsInstalled() should return false after uninstall")
	}
}

func TestInstallPromptToClaudeCode(t *testing.T) {
	tmpDir := t.TempDir()

	agent := Agent{
		Name:        "Claude Code Test",
		Type:        ClaudeCode,
		ConfigPath:  tmpDir,
		PromptPath:  filepath.Join(tmpDir, "CLAUDE.md"),
		Description: "Test agent",
	}

	testPrompt := "# LazyMentor - System Prompt\nTest content"

	err := agent.InstallPrompt(testPrompt)
	if err != nil {
		t.Fatalf("InstallPrompt() error = %v", err)
	}

	// Verify file was created
	installedPath := filepath.Join(tmpDir, "CLAUDE.md")
	content, err := os.ReadFile(installedPath)
	if err != nil {
		t.Fatalf("Failed to read installed file: %v", err)
	}

	if string(content) != testPrompt {
		t.Errorf("Installed content = %q, want %q", string(content), testPrompt)
	}
}

func TestInstallPromptMergesForClaudeCode(t *testing.T) {
	tmpDir := t.TempDir()

	agent := Agent{
		Name:        "Claude Code Test",
		Type:        ClaudeCode,
		ConfigPath:  tmpDir,
		PromptPath:  filepath.Join(tmpDir, "CLAUDE.md"),
		Description: "Test agent",
	}

	// Install first with existing content
	existingContent := "# Project Notes\n\nSome notes here"
	if err := os.WriteFile(agent.PromptPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing file: %v", err)
	}

	testPrompt := "# LazyMentor - System Prompt\nTest content"

	err := agent.InstallPrompt(testPrompt)
	if err != nil {
		t.Fatalf("InstallPrompt() error = %v", err)
	}

	// Verify file was created and contains both
	content, err := os.ReadFile(agent.PromptPath)
	if err != nil {
		t.Fatalf("Failed to read installed file: %v", err)
	}

	// Should contain existing content
	if !contains(string(content), "# Project Notes") {
		t.Error("Installed content should contain original content")
	}

	// Should contain lazymentor
	if !contains(string(content), LazyMentorMarker) {
		t.Error("Installed content should contain lazymentor")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
