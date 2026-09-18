package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// TestPromptTotality proves the dock renders every engine prompt kind (ADR 0045).
// promptControls is the seam controls dispatches through, so a new Chooser
// capability — a new engine.PromptKind — that no case here handles fails this
// test until its rendering is added.
func TestPromptTotality(t *testing.T) {
	c := newClient(t)
	for _, kind := range engine.PromptKinds() {
		if _, covered := c.g.promptControls(kind); !covered {
			t.Errorf("prompt kind %v has no dock rendering; add a case to promptControls", kind)
		}
	}
}

// TestOptionKindTotality proves the option chooser draws every option-widget
// kind (ADR 0045), so a new option shape cannot fall through to a silent
// catch-all: it must be a named optionKind with its own case in optionControls.
func TestOptionKindTotality(t *testing.T) {
	c := newClient(t)
	for _, kind := range optionKinds() {
		if _, covered := c.g.optionControls(kind); !covered {
			t.Errorf("option kind %d has no widget; add a case to optionControls", kind)
		}
	}
}
