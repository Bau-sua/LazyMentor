package config

import (
	"runtime"
	"testing"
)

func TestDetectOS(t *testing.T) {
	os := DetectOS()

	switch runtime.GOOS {
	case "linux":
		if os != Linux {
			t.Errorf("DetectOS() = %s, want %s", os, Linux)
		}
	case "darwin":
		if os != MacOS {
			t.Errorf("DetectOS() = %s, want %s", os, MacOS)
		}
	case "windows":
		if os != Windows {
			t.Errorf("DetectOS() = %s, want %s", os, Windows)
		}
	}
}

func TestGetOSInfo(t *testing.T) {
	info := GetOSInfo()

	if info.OS == "" {
		t.Error("OS should not be empty")
	}

	if info.Name == "" {
		t.Error("Name should not be empty")
	}

	if info.Arch == "" {
		t.Error("Arch should not be empty")
	}

	if info.GoOS != runtime.GOOS {
		t.Errorf("GoOS = %s, want %s", info.GoOS, runtime.GOOS)
	}

	if info.GoArch != runtime.GOARCH {
		t.Errorf("GoArch = %s, want %s", info.GoArch, runtime.GOARCH)
	}
}

func TestFormatOS(t *testing.T) {
	formatted := FormatOS()

	if formatted == "" {
		t.Error("FormatOS() should not return empty string")
	}

	// Should contain OS name
	info := GetOSInfo()
	if !contains(formatted, info.Name) {
		t.Errorf("FormatOS() = %s, should contain OS name %s", formatted, info.Name)
	}

	// Should contain GOOS
	if !contains(formatted, runtime.GOOS) {
		t.Errorf("FormatOS() = %s, should contain %s", formatted, runtime.GOOS)
	}

	// Should contain GOARCH
	if !contains(formatted, runtime.GOARCH) {
		t.Errorf("FormatOS() = %s, should contain %s", formatted, runtime.GOARCH)
	}
}

func TestGetOSName(t *testing.T) {
	tests := []struct {
		os       OS
		expected string
	}{
		{Linux, "Linux"},
		{MacOS, "macOS"},
		{Windows, "Windows"},
		{OS("freebsd"), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.os), func(t *testing.T) {
			result := getOSName(tt.os)
			if result != tt.expected {
				t.Errorf("getOSName(%s) = %s, want %s", tt.os, result, tt.expected)
			}
		})
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
