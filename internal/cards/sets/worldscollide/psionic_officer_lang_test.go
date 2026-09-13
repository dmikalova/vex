package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Psionic Officer Lang
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	After an enemy Creature reaps, archive the top card of your deck.
func TestPsionicOfficerLang(t *testing.T) {
	t.Run("archives the top card of your deck after an enemy creature reaps", func(t *testing.T) {
		var topCard, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(PsionicOfficerLang),
				Hand:   ct.DeckOf(card.House.StarAlliance, 6), // full hand: no end-of-turn draw
				Deck:   ct.Cards(ct.Bind(&topCard, ct.Creature())),
			},
			P2: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Untamed)
		h.P2.Reap(foe)

		h.Expect(topCard).At(ct.Archives)
	})
}
