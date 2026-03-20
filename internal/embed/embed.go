// Package embed contains the embedded lazymentor.md prompt
package embed

import (
	_ "embed"
)

//go:embed lazymentor.md
var LazyMentorPrompt string // lazymentor.md is embedded at build time
