// Package config provides platform and OS detection utilities
package config

import (
	"fmt"
	"runtime"
)

// OS represents a supported operating system
type OS string

const (
	Linux   OS = "linux"
	MacOS   OS = "darwin"
	Windows OS = "windows"
)

// DetectOS returns the current operating system
func DetectOS() OS {
	switch runtime.GOOS {
	case "linux":
		return Linux
	case "darwin":
		return MacOS
	case "windows":
		return Windows
	default:
		return Linux
	}
}

// OSInfo contains information about the operating system
type OSInfo struct {
	OS     OS
	Name   string
	Arch   string
	GoOS   string
	GoArch string
}

// GetOSInfo returns detailed OS information
func GetOSInfo() OSInfo {
	return OSInfo{
		OS:     DetectOS(),
		Name:   getOSName(DetectOS()),
		Arch:   runtime.GOARCH,
		GoOS:   runtime.GOOS,
		GoArch: runtime.GOARCH,
	}
}

func getOSName(os OS) string {
	switch os {
	case Linux:
		return "Linux"
	case MacOS:
		return "macOS"
	case Windows:
		return "Windows"
	default:
		return "Unknown"
	}
}

// FormatOS returns a formatted string for the current OS
func FormatOS() string {
	info := GetOSInfo()
	return fmt.Sprintf("%s (%s/%s)", info.Name, info.GoOS, info.GoArch)
}
