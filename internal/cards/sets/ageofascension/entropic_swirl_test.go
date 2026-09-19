package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Entropic Swirl
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature. For each trait that creature has, deal 2 damage to the chosen creature. Gain 1 Æmber.
func TestEntropicSwirl(t *testing.T) {
	var foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			Hand:  ct.Cards(EntropicSwirl),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(
				ct.Power(6),
				ct.Traits(engine.Beast, engine.Mutant),
			))),
		},
	})

	h.P1.Play(EntropicSwirl)

	// Two traits: 4 damage dealt, 2 Æmber gained.
	h.Expect(foe).Damage(4)
	h.P1.ExpectAmber(2)
}
