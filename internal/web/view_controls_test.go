package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
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

// TestPromptSourceReadsTheAskingCopysHouse pins the prompt source face to the
// card that actually raised the prompt. With two copies of one card in play, the
// face must show the asking copy's house — not the house of whichever copy the
// battleline happens to hold first.
func TestPromptSourceReadsTheAskingCopysHouse(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	board := c.board()
	if len(board) != 2 {
		t.Fatalf("the battleline holds %d creatures, want 2 copies", len(board))
	}
	// The copy nearer the left flank is resolving as another house; the prompt is
	// raised by the other copy, which still belongs to its printed house.
	c.g.eng().SetLastingHouse(board[0], engine.Dis)

	answer := c.ask(board[1], "Choose a creature", false, board)
	c.await("the prompt to go up", c.g.choosing)

	// app.If only draws inside a parent element, so the block is rendered in one.
	html := app.HTMLString(app.Div().Body(c.g.promptSourceHeader()))
	if !strings.Contains(html, houseClasses(testHouse)) {
		t.Errorf("the prompt source face is not drawn in %v: %s", testHouse, html)
	}
	if strings.Contains(html, houseClasses(engine.Dis)) {
		t.Errorf("the prompt source face took the other copy's house: %s", html)
	}
	if strings.Contains(html, "icon-house--changed") {
		t.Errorf("the prompt source is marked as a changed house, but the asking "+
			"copy belongs to its printed house: %s", html)
	}

	c.g.chooseCandidate(c.ctx, board[1])
	c.settle()
	if !answer.ok || answer.id != board[1] {
		t.Errorf("the effect was answered %v, want card %d", answer, board[1])
	}
	c.await("the prompt to come down", func() bool { return !c.g.choosing() })
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
