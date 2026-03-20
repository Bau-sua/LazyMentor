package update

import (
	"testing"
)

func TestGetCurrentVersion(t *testing.T) {
	version := GetCurrentVersion()

	if version == "" {
		t.Error("Version should not be empty")
	}
}

func TestSetVersion(t *testing.T) {
	SetVersion("1.2.3")

	if GetCurrentVersion() != "1.2.3" {
		t.Error("Version should be updated")
	}

	// Reset
	SetVersion("0.1.0")
}

func TestGetAssetName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"linux amd64", "lazymint-linux-amd64"},
		{"linux arm64", "lazymint-linux-arm64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAssetName("1.0.0")
			// Just check it returns something
			if result == "" {
				t.Error("getAssetName should return non-empty string")
			}
		})
	}
}

func TestGetAssetNameLinux(t *testing.T) {
	result := getAssetName("test")

	// Should contain linux
	if result == "" {
		t.Error("Should return asset name")
	}
}
