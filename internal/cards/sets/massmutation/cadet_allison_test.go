package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cadet Allison
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  8
//	Traits: Human
//
//	Play/Reap: Discard a random card from your hand -> its house becomes your active house.
func TestCadetAllison(t *testing.T) {
	var allison, spare ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(ct.Bind(&allison, CadetAllison)),
			Hand:   ct.Cards(ct.Bind(&spare, ct.Tactic(ct.OfHouse(card.House.Untamed)))),
		},
	})

	h.P1.Reap(allison)

	// Reaping gains 1 Æmber and discards the random hand card.
	h.P1.ExpectAmber(1)
	h.Expect(spare).At(ct.Discard)
}
