package anomalyexpansion

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Memolith
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Action: Choose one:
//	- Put a tactic card from your hand faceup under Memolith
//	- Trigger the play effect of a Tactic grafted onto Memolith.
func TestMemolith(t *testing.T) {
	// Memolith's action grafts a Tactic from hand; readied and used again, it
	// triggers that grafted Tactic's play effect while it stays grafted.
	var memo, refrain, dummy ct.Card
	// Warriors' Refrain (Play: stun each creature with power 3 or lower) stands in
	// for a real Tactic to graft; it is built here rather than imported from its
	// own set so this test stays self-contained.
	refrainDef := card.Build(
		"Warriors' Refrain",
		card.House.Brobnar,
		card.Type.Tactic,
		card.Rarity.Common,
		card.WithBonus(card.Bonus.Aember),
		card.WithAbility(
			card.Trigger.Play, card.Stun{
				Target: card.Target.EachCreature.PowerAtMost(3),
			}),
	)
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			InPlay: ct.Cards(
				ct.Bind(&memo, Memolith),
				ct.Bind(&dummy, ct.Creature(ct.Power(3))),
			),
			Hand: ct.Cards(ct.Bind(&refrain, refrainDef)),
		},
	})

	h.P1.UseAction(Memolith)
	h.P1.ClickOption("put a Tactic")
	h.Expect(refrain).At(ct.Under)

	// Ready Memolith to use it again for the other half of the choice.
	memo.Ready()
	h.P1.UseAction(Memolith)
	h.P1.ClickOption("trigger the play effect")

	// Warriors' Refrain's play effect stunned the low-power friendly creature,
	// and the Tactic stayed grafted under Memolith.
	h.Expect(dummy).Stunned(true)
	h.Expect(refrain).At(ct.Under)
}
