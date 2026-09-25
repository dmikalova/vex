package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Horizon Saber
//
//	House:  Logos
//	Type:   Gigantic Creature
//	Rarity: Special
//	Power:  11
//	Armor:  2
//	Traits: Robot
//
//	Play/Fight/Reap: Search your deck and discard pile for a card, reveal it, and put it into your archives. Shuffle your discard pile into your deck.
func TestHorizonSaber(t *testing.T) {
	var saber, target, buried ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:   card.House.Logos,
			InPlay:  ct.Cards(ct.Bind(&saber, HorizonSaber)),
			Deck:    ct.Cards(ct.Bind(&target, ct.Creature())),
			Discard: ct.Cards(ct.Bind(&buried, ct.Creature())),
		},
	})

	h.P1.Reap(saber)
	h.P1.ClickCard(target) // archive the chosen deck card

	h.Expect(target).At(ct.Archives) // the chosen card goes to archives
	h.Expect(buried).At(ct.Deck)     // the discard pile is shuffled into the deck
}
