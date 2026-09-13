package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Memolith
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Action: Choose one:
//	- Put a Tactic card from your hand faceup under Memolith
//	- Trigger the play effect of a Tactic grafted onto Memolith.
func TestMemolith(t *testing.T) {
	// Memolith's action grafts an action from hand; readied and used again, it
	// triggers that grafted action's play effect while it stays grafted.
	var memo, brew, dummy ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			InPlay: ct.Cards(
				ct.Bind(&memo, Memolith),
				ct.Bind(&dummy, ct.Creature(ct.Power(3))),
			),
			Hand: ct.Cards(ct.Bind(&brew, ChieftainsBrew)),
		},
	})

	h.P1.UseAction(Memolith)
	h.P1.ClickOption("put a Tactic")
	h.Expect(brew).At(ct.Under)

	// Ready Memolith to use it again for the other half of the choice.
	memo.Ready()
	h.P1.UseAction(Memolith)
	h.P1.ClickOption("trigger the play effect")

	// Chieftain's Brew's play effect gave the friendly creature two +1 power
	// counters, and the action stayed grafted under Memolith.
	h.Expect(dummy).Power(5)
	h.Expect(brew).At(ct.Under)
}
