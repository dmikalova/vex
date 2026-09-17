package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// TestCardFaceBonusAndEnhance checks a card's printed bonus icons render as the
// left-edge badge strip and an Enhance source's "Enhance <icons>" line renders as
// an ordinary final line of the rules text, with glyphs rather than icon words.
func TestCardFaceBonusAndEnhance(t *testing.T) {
	def := engine.NewCard("Splinterish", engine.Shadows, engine.Creature, engine.Common,
		engine.WithPower(1),
		engine.WithBonus(engine.BonusAember),
		engine.WithEnhance(engine.BonusDamage, engine.BonusDamage))
	html := app.HTMLString(printedFace(&def).Render())

	for _, want := range []string{"card-bonuses", "aember.svg", "card-rules", "Enhance", "damage.svg"} {
		if !strings.Contains(html, want) {
			t.Errorf("card face missing %q", want)
		}
	}
	// the Enhance line shows the damage glyph, not the word "Damage".
	if strings.Contains(html, "Enhance Damage") {
		t.Errorf("Enhance line should render icons, not the word Damage: %s", html)
	}
	// the line is inlined in the rules text box, not a text box of its own.
	if strings.Contains(html, "card-enhance") {
		t.Errorf("Enhance line should be part of card-rules, not its own div: %s", html)
	}
}
