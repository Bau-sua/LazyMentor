package nvim

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	cfg := Detect()

	// Should detect Neovim on this system
	if !cfg.Installed {
		t.Skip("Neovim not installed, skipping test")
	}

	if cfg.Version == "" {
		t.Error("Version should not be empty when installed")
	}

	if cfg.ConfigDir == "" {
		t.Error("ConfigDir should not be empty")
	}
}

func TestDetectLeaderKey(t *testing.T) {
	// Test with temp directory
	tmpDir := t.TempDir()

	// Create a mock init.lua with leader key
	initLua := `vim.g.mapleader = ","`
	initPath := filepath.Join(tmpDir, "init.lua")
	os.WriteFile(initPath, []byte(initLua), 0644)

	leader := detectLeaderKey(tmpDir)

	// Should detect the leader
	if leader == "" {
		t.Error("Should detect leader key")
	}
}

func TestIsLazyVim(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with lazy-lock.json
	lockFile := filepath.Join(tmpDir, "lazy-lock.json")
	os.WriteFile(lockFile, []byte("{}"), 0644)

	if !isLazyVim(tmpDir) {
		t.Error("Should detect LazyVim from lazy-lock.json")
	}
}

func TestDetectPlugins(t *testing.T) {
	tmpDir := t.TempDir()

	// Create lua directory with telescope plugin
	luaDir := filepath.Join(tmpDir, "lua", "telescope.nvim")
	os.MkdirAll(luaDir, 0755)

	plugins := detectPlugins(tmpDir)

	// Should find telescope plugin
	found := false
	for _, p := range plugins {
		if p == "telescope.nvim" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Should detect telescope.nvim plugin")
	}
}

func TestConfigFormatForPrompt(t *testing.T) {
	cfg := Config{
		Installed: true,
		Version:   "0.10.0",
		IsLazyVim: true,
		ConfigDir: "/home/user/.config/nvim",
		LeaderKey: "Space",
		Plugins:   []string{"telescope.nvim", "harpoon"},
	}

	output := cfg.FormatForPrompt()

	if output == "" {
		t.Error("FormatForPrompt should return non-empty string")
	}

	if !contains(output, "Neovim: 0.10.0") {
		t.Error("Should contain version")
	}

	if !contains(output, "LazyVim") {
		t.Error("Should contain LazyVim indicator")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
