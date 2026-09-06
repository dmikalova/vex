package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Destroy Them All!
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy an artifact and a creature and an upgrade.
func TestDestroyThemAll(t *testing.T) {
	var relic, victim, host, boon ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Mars,
			Hand:  ct.Cards(DestroyThemAll),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&relic, ct.Artifact()),
				ct.Bind(&victim, ct.Creature(ct.Power(3))),
				ct.Upgraded(
					ct.Bind(&host, ct.Creature(ct.Power(5))),
					ct.Bind(&boon, ct.Upgrade(ct.PowerBonus(2))),
				),
			),
		},
	})

	h.P1.Play(DestroyThemAll)
	// The lone artifact and lone upgrade are forced; the creature is chosen.
	h.P1.ClickCard(victim)

	h.Expect(relic).At(ct.Discard)
	h.Expect(victim).At(ct.Discard)
	h.Expect(boon).At(ct.Discard)
	h.Expect(host).At(ct.PlayArea)
}
