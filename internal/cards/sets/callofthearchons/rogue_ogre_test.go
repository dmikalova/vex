package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Rogue Ogre
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant • Mutant
//
//	At the end of your turn, if you played exactly 1 card this turn, heal 2 damage from Rogue Ogre, and Rogue Ogre captures 1 Æmber from your opponent.
func TestRogueOgre(t *testing.T) {
	filler := ct.Tactic(ct.OfHouse(card.House.Brobnar))
	other := ct.Tactic(ct.OfHouse(card.House.Brobnar))
	setup := func(ogre *ct.Card) ct.Setup {
		return ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(filler, other),
				InPlay: ct.Cards(ct.Bind(ogre, RogueOgre)),
			},
			P2: ct.Side{Amber: 4},
		}
	}

	t.Run("heals and captures after a one-card turn", func(t *testing.T) {
		var ogre ct.Card
		h := ct.Play(t, setup(&ogre))
		ogre.Damaged(3)

		h.P1.Play(filler)
		h.P1.EndTurn()

		h.Expect(ogre).Damage(1).AmberOn(1)
		h.P2.ExpectAmber(3)
	})

	t.Run("does nothing after a two-card turn", func(t *testing.T) {
		var ogre ct.Card
		h := ct.Play(t, setup(&ogre))
		ogre.Damaged(3)

		h.P1.Play(filler)
		h.P1.Play(other)
		h.P1.EndTurn()

		h.Expect(ogre).Damage(3).AmberOn(0)
		h.P2.ExpectAmber(4)
	})
}
