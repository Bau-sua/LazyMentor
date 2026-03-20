// Package tui provides the terminal user interface using bubbletea
package tui

import (
	"github.com/charmbracelet/bubbletea"
)

// Run starts the TUI application
func Run() error {
	model := NewWelcomeModel()
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
