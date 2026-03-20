// Package nvim provides Neovim detection and configuration discovery
package nvim

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Config holds information about the user's Neovim configuration
type Config struct {
	Installed bool     // Is Neovim installed?
	Version   string   // Neovim version
	IsLazyVim bool     // Is using LazyVim?
	IsNeovim  bool     // Is using regular Neovim (not LazyVim)
	ConfigDir string   // Config directory path
	LeaderKey string   // Map leader key
	Plugins   []string // Detected plugins
}

// Detect finds and analyzes the Neovim configuration
func Detect() Config {
	cfg := Config{
		Installed: false,
	}

	// Check if nvim command exists
	if !isNVimInstalled() {
		return cfg
	}
	cfg.Installed = true

	// Get version
	cfg.Version = getNVimVersion()

	// Find config directory
	cfg.ConfigDir = getConfigDir()

	// Check if LazyVim
	cfg.IsLazyVim = isLazyVim(cfg.ConfigDir)

	// Get plugins
	cfg.Plugins = detectPlugins(cfg.ConfigDir)

	// Get leader key
	cfg.LeaderKey = detectLeaderKey(cfg.ConfigDir)

	return cfg
}

func isNVimInstalled() bool {
	_, err := exec.LookPath("nvim")
	return err == nil
}

func getNVimVersion() string {
	cmd := exec.Command("nvim", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	if scanner.Scan() {
		line := scanner.Text()
		// First line is typically "NVIM v0.10.0" or "NVIM v0.10.1"
		parts := strings.Split(line, " ")
		if len(parts) >= 2 {
			return strings.TrimPrefix(parts[1], "v")
		}
	}
	return "unknown"
}

func getConfigDir() string {
	home := os.Getenv("HOME")
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "nvim-data")
	default:
		// Check XDG_CONFIG_HOME first
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, "nvim")
		}
		return filepath.Join(home, ".config", "nvim")
	}
}

func isLazyVim(configDir string) bool {
	// Check for lazy.nvim in lua directory
	pluginFile := filepath.Join(configDir, "lua", "lazy.nvim")

	if _, err := os.Stat(pluginFile); err == nil {
		return true
	}

	// Check lazy-lock.json (LazyVim creates this)
	lockFile := filepath.Join(configDir, "lazy-lock.json")
	if _, err := os.Stat(lockFile); err == nil {
		return true
	}

	// Check for lazy.lua in config
	initLua := filepath.Join(configDir, "init.lua")
	if data, err := os.ReadFile(initLua); err == nil {
		content := string(data)
		if strings.Contains(content, "lazy.nvim") || strings.Contains(content, "LazyVim") {
			return true
		}
	}

	return false
}

func detectPlugins(configDir string) []string {
	var plugins []string

	// Common plugin indicators
	pluginIndicators := map[string][]string{
		"telescope.nvim":      {"telescope", "nvim-telescope"},
		"harpoon":             {"harpoon", "harpoon2"},
		"neo-tree.nvim":       {"neo-tree", "neotree"},
		"oil.nvim":            {"oil", "oil.nvim"},
		"toggleterm.nvim":     {"toggleterm", "toggleterm"},
		"fzf":                 {"fzf", "telescope-fzf"},
		"fzf.lua":             {"fzf-lua", "telescope-fzf-native-native"},
		"diffview.nvim":       {"diffview", "diffview"},
		"gitsigns.nvim":       {"gitsigns", "git-signs"},
		"lspsaga.nvim":        {"lspsaga", "Saga"},
		"telescope-ui-select": {"telescope", "ui-select"},
	}

	// Search in lua directory for plugins
	dirsToSearch := []string{
		filepath.Join(configDir, "lua"),
		filepath.Join(configDir, "plugin"),
		configDir,
	}

	// Also check lazy.nvim plugin cache
	pluginCache := filepath.Join(configDir, "lazy-lock.json")
	if _, err := os.Stat(pluginCache); err == nil {
		plugins = append(plugins, "lazy.nvim (LazyVim detected)")
	}

	for _, dir := range dirsToSearch {
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		entries, _ := os.ReadDir(dir)
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()

			for pluginName, indicators := range pluginIndicators {
				for _, indicator := range indicators {
					if strings.Contains(strings.ToLower(name), indicator) {
						plugins = append(plugins, pluginName)
						break
					}
				}
			}
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	unique := []string{}
	for _, p := range plugins {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}

	return unique
}

func detectLeaderKey(configDir string) string {
	// Look for mapleader in lua or init.lua
	filesToCheck := []string{
		filepath.Join(configDir, "init.lua"),
		filepath.Join(configDir, "lua", "core", "options.lua"),
		filepath.Join(configDir, "lua", "options.lua"),
		filepath.Join(configDir, "lua", "core", "keymaps.lua"),
		filepath.Join(configDir, "lua", "keymaps.lua"),
	}

	for _, file := range filesToCheck {
		if data, err := os.ReadFile(file); err == nil {
			// Simple line-by-line pattern matching
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.Contains(line, "mapleader") && strings.Contains(line, "=") {
					parts := strings.Split(line, "=")
					if len(parts) >= 2 {
						right := strings.TrimSpace(parts[1])
						right = strings.Trim(right, " \t,\"'\");")
						if len(right) > 0 && len(right) <= 3 {
							return right
						}
					}
				}
			}
		}
	}

	return "Space" // Default leader key
}

// FormatForPrompt returns a formatted string for use in the prompt
func (c Config) FormatForPrompt() string {
	if !c.Installed {
		return "Neovim: Not detected\n\nTip: LazyMentor works best when you have Neovim open!"
	}

	var lines []string

	lines = append(lines, "## User Environment")
	lines = append(lines, "- Neovim: "+c.Version+" ✓")
	lines = append(lines, "- Config: "+c.ConfigDir)

	if c.IsLazyVim {
		lines = append(lines, "- Distro: LazyVim ✓")
	} else {
		lines = append(lines, "- Distro: Neovim (not LazyVim)")
	}

	if c.LeaderKey != "" {
		lines = append(lines, "- Leader key: "+c.LeaderKey)
	}

	if len(c.Plugins) > 0 {
		lines = append(lines, "- Plugins: "+strings.Join(c.Plugins, ", "))
	} else {
		lines = append(lines, "- Plugins: None detected")
	}

	return strings.Join(lines, "\n")
}
