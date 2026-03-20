package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lazymentor/lazymint/internal/agents"
	"github.com/lazymentor/lazymint/internal/config"
	"github.com/lazymentor/lazymint/internal/embed"
)

// Screen represents the current screen
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenPreflight
	ScreenAgentSelect
	ScreenInstalling
	ScreenSuccess
	ScreenError
)

// InstallMsg signals that installation completed
type InstallMsg struct {
	Err error
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#6366F1")).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A1A1AA")).
			MarginTop(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#71717A"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)

	checkboxChecked   = lipgloss.Style{}.Foreground(lipgloss.Color("#22C55E")).Render("✓")
	checkboxUnchecked = lipgloss.Style{}.Foreground(lipgloss.Color("#71717A")).Render("○")
)

// Model represents the TUI state
type Model struct {
	Screen            Screen
	SelectedAgent     int
	Agents            []agents.Agent
	PreflightNVim     bool
	PreflightLazyVim  bool
	InstallingMessage string
	ErrorMessage      string
}

// NewWelcomeModel creates a new welcome model
func NewWelcomeModel() Model {
	return Model{
		Screen:           ScreenWelcome,
		Agents:           agents.DetectAgents(),
		PreflightNVim:    true,
		PreflightLazyVim: true,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case InstallMsg:
		if msg.Err != nil {
			m.Screen = ScreenError
			m.ErrorMessage = msg.Err.Error()
		} else {
			m.Screen = ScreenSuccess
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			switch m.Screen {
			case ScreenWelcome:
				m.Screen = ScreenPreflight
			case ScreenPreflight:
				m.Screen = ScreenAgentSelect
			case ScreenAgentSelect:
				if len(m.Agents) > 0 {
					m.Screen = ScreenInstalling
					return m, installPromptCmd(m.SelectedAgent, m.Agents[m.SelectedAgent])
				}
			}
		case "up", "k":
			if m.Screen == ScreenAgentSelect && m.SelectedAgent > 0 {
				m.SelectedAgent--
			}
		case "down", "j":
			if m.Screen == ScreenAgentSelect && m.SelectedAgent < len(m.Agents)-1 {
				m.SelectedAgent++
			}
		}
	}
	return m, nil
}

// View renders the model
func (m Model) View() string {
	var s strings.Builder

	switch m.Screen {
	case ScreenWelcome:
		s.WriteString(welcomeView(m))
	case ScreenPreflight:
		s.WriteString(preflightView(m))
	case ScreenAgentSelect:
		s.WriteString(agentSelectView(m))
	case ScreenInstalling:
		s.WriteString(installingView(m))
	case ScreenSuccess:
		s.WriteString(successView(m))
	case ScreenError:
		s.WriteString(errorView(m))
	}

	s.WriteString("\n")
	s.WriteString(normalStyle.Render("  Press "))
	s.WriteString(selectedStyle.Render("q"))
	s.WriteString(normalStyle.Render(" to quit, "))
	s.WriteString(selectedStyle.Render("enter"))
	s.WriteString(normalStyle.Render(" to continue\n"))

	return s.String()
}

func welcomeView(m Model) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" LazyMentor Installer "))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render("Your LazyVim learning companion"))
	s.WriteString("\n\n")
	s.WriteString(headerStyle.Render(" Welcome "))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  %s %s\n", normalStyle.Render("OS:"), config.FormatOS()))
	s.WriteString("\n")
	s.WriteString(normalStyle.Render("  This installer will help you set up LazyMentor,"))
	s.WriteString("\n")
	s.WriteString(normalStyle.Render("  a mentor that teaches you LazyVim keybindings."))
	s.WriteString("\n\n")
	s.WriteString(selectedStyle.Render("  ▶ Press Enter to continue..."))
	return s.String()
}

func preflightView(m Model) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" LazyMentor Installer "))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render("Your LazyVim learning companion"))
	s.WriteString("\n\n")
	s.WriteString(headerStyle.Render(" Pre-flight Checklist "))
	s.WriteString("\n\n")
	s.WriteString(normalStyle.Render("  This tool works best when you practice along."))
	s.WriteString("\n")
	s.WriteString(normalStyle.Render("  Make sure you have Neovim open!"))
	s.WriteString("\n\n")

	// Checkbox 1
	cb1 := checkboxUnchecked
	if m.PreflightNVim {
		cb1 = checkboxChecked
	}
	s.WriteString(fmt.Sprintf("  %s %s\n", cb1, normalStyle.Render("I have Neovim open and ready")))
	s.WriteString("\n")
	s.WriteString(normalStyle.Render("  (You can skip this, but you'll learn faster with it open!)"))
	s.WriteString("\n\n")
	s.WriteString(selectedStyle.Render("  ▶ Press Enter to continue..."))
	return s.String()
}

func agentSelectView(m Model) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" LazyMentor Installer "))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render("Your LazyVim learning companion"))
	s.WriteString("\n\n")
	s.WriteString(headerStyle.Render(" Select Agent "))
	s.WriteString("\n\n")

	if len(m.Agents) == 0 {
		s.WriteString(errorStyle.Render("  ✗ No supported agents detected"))
		s.WriteString("\n\n")
		s.WriteString(normalStyle.Render("  Install OpenCode or Claude Code to continue."))
		s.WriteString("\n")
		s.WriteString(normalStyle.Render("  Or use manual installation mode."))
	} else {
		s.WriteString(normalStyle.Render("  Choose where to install LazyMentor:"))
		s.WriteString("\n\n")
		for i, agent := range m.Agents {
			prefix := "  "
			if i == m.SelectedAgent {
				prefix = "▶ "
			}
			style := normalStyle
			if i == m.SelectedAgent {
				style = selectedStyle
			}
			s.WriteString(fmt.Sprintf("%s%s%s %s\n", prefix, style.Render(agent.Name), normalStyle.Render(" - "+agent.Description), normalStyle.Render("("+agent.ConfigPath+")")))
		}
		s.WriteString("\n")
		s.WriteString(normalStyle.Render("  Use "))
		s.WriteString(selectedStyle.Render("↑/↓"))
		s.WriteString(normalStyle.Render(" or "))
		s.WriteString(selectedStyle.Render("j/k"))
		s.WriteString(normalStyle.Render(" to navigate, "))
		s.WriteString(selectedStyle.Render("Enter"))
		s.WriteString(normalStyle.Render(" to install"))
	}

	return s.String()
}

func installingView(m Model) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" LazyMentor Installer "))
	s.WriteString("\n\n")
	s.WriteString(headerStyle.Render(" Installing "))
	s.WriteString("\n\n")
	s.WriteString(normalStyle.Render("  " + m.InstallingMessage))
	s.WriteString("\n\n")
	return s.String()
}

func successView(m Model) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" LazyMentor Installer "))
	s.WriteString("\n\n")
	s.WriteString(successStyle.Render("  ✓ Installation Complete!"))
	s.WriteString("\n\n")
	s.WriteString(normalStyle.Render("  LazyMentor has been installed successfully."))
	s.WriteString("\n")
	s.WriteString(normalStyle.Render("  Restart your agent to start learning LazyVim!"))
	s.WriteString("\n\n")
	return s.String()
}

func errorView(m Model) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" LazyMentor Installer "))
	s.WriteString("\n\n")
	s.WriteString(errorStyle.Render("  ✗ Error"))
	s.WriteString("\n\n")
	s.WriteString(normalStyle.Render("  " + m.ErrorMessage))
	s.WriteString("\n\n")
	return s.String()
}

// installPrompt returns a command that installs the prompt asynchronously
func installPromptCmd(selectedAgent int, agent agents.Agent) tea.Cmd {
	return func() tea.Msg {
		prompt := loadLocalPrompt()
		if prompt == "" {
			prompt = embed.LazyMentorPrompt
		}

		if err := agent.InstallPrompt(prompt); err != nil {
			return InstallMsg{Err: err}
		}

		return InstallMsg{Err: nil}
	}
}

func loadLocalPrompt() string {
	data, err := os.ReadFile("lazymentor.md")
	if err != nil {
		return ""
	}
	return string(data)
}
