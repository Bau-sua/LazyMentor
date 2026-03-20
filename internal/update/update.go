// Package update provides automatic update checking and installation
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Version represents a release version
type Version struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}

// CheckResult contains the result of checking for updates
type CheckResult struct {
	CurrentVersion string
	LatestVersion  string
	UpdateNeeded   bool
	DownloadURL    string
	AssetName      string
	AssetSize      int64
	ReleaseURL     string
}

// GetCurrentVersion returns the current version from the binary
func GetCurrentVersion() string {
	// This is set at compile time via -ldflags
	return version
}

// CheckForUpdates checks if there's a newer version available
func CheckForUpdates(repo string) (*CheckResult, error) {
	currentVer := GetCurrentVersion()

	// Fetch latest release from GitHub
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add User-Agent to avoid rate limiting
	req.Header.Set("User-Agent", "LazyMentor-Installer")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to check updates: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release Version
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	latestVer := strings.TrimPrefix(release.TagName, "v")

	// Find the correct asset for this platform
	assetName := getAssetName(latestVer)
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, release.TagName, assetName)

	result := &CheckResult{
		CurrentVersion: currentVer,
		LatestVersion:  latestVer,
		UpdateNeeded:   latestVer != currentVer,
		DownloadURL:    downloadURL,
		AssetName:      assetName,
		ReleaseURL:     release.HTMLURL,
	}

	return result, nil
}

// DownloadAndReplace downloads the new version and replaces the current binary
func DownloadAndReplace(url string, currentBinaryPath string) error {
	// Create temp file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "lazymint-update")

	// Download
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Make executable
	if err := os.Chmod(tempFile, 0755); err != nil {
		return fmt.Errorf("failed to make executable: %w", err)
	}

	// Create backup
	backupPath := currentBinaryPath + ".backup"
	if err := os.Rename(currentBinaryPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Move new binary to final location
	if err := os.Rename(tempFile, currentBinaryPath); err != nil {
		// Try to restore backup
		os.Rename(backupPath, currentBinaryPath)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Remove backup on success
	os.Remove(backupPath)

	return nil
}

// getAssetName returns the correct binary name for the current platform
func getAssetName(version string) string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "linux":
		switch arch {
		case "amd64":
			return fmt.Sprintf("lazymint-linux-amd64")
		case "arm64":
			return fmt.Sprintf("lazymint-linux-arm64")
		}
	case "darwin":
		switch arch {
		case "amd64":
			return fmt.Sprintf("lazymint-darwin-amd64")
		case "arm64":
			return fmt.Sprintf("lazymint-darwin-arm64")
		}
	case "windows":
		return fmt.Sprintf("lazymint-windows-amd64.exe")
	}

	// Fallback
	return fmt.Sprintf("lazymint-%s-%s", os, arch)
}

// version is set at compile time
var version = "0.1.0"

// SetVersion sets the version (called from main via ldflags)
func SetVersion(v string) {
	version = v
}
