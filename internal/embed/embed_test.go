package embed

import (
	"strings"
	"testing"
)

func TestLazyMentorPromptLoaded(t *testing.T) {
	if LazyMentorPrompt == "" {
		t.Error("LazyMentorPrompt should not be empty")
	}
}

func TestLazyMentorPromptContainsMarker(t *testing.T) {
	expected := "# LazyMentor - System Prompt"

	if !strings.Contains(LazyMentorPrompt, expected) {
		t.Errorf("LazyMentorPrompt should contain %q", expected)
	}
}

func TestLazyMentorPromptContainsRules(t *testing.T) {
	rules := []string{
		"**NUNCA** generes",
		"tablas Markdown",
		"Respuestas",
	}

	for _, rule := range rules {
		if !strings.Contains(LazyMentorPrompt, rule) {
			t.Errorf("LazyMentorPrompt should contain %q", rule)
		}
	}
}

func TestLazyMentorPromptContainsReglaDeOro(t *testing.T) {
	if !strings.Contains(LazyMentorPrompt, "Regla de Oro") {
		t.Error("LazyMentorPrompt should contain 'Regla de Oro'")
	}
}

func TestLazyMentorPromptContainsExamples(t *testing.T) {
	if !strings.Contains(LazyMentorPrompt, "cómo creo un archivo nuevo") {
		t.Error("LazyMentorPrompt should contain example about creating files")
	}
}

func TestLazyMentorPromptMinimumLength(t *testing.T) {
	minLength := 500

	if len(LazyMentorPrompt) < minLength {
		t.Errorf("LazyMentorPrompt length = %d, want at least %d", len(LazyMentorPrompt), minLength)
	}
}
