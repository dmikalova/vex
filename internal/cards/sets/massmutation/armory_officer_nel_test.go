package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Armory Officer Nel
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Alien
//
//	After an upgrade enters play, draw a card.
//	Enhance Draw.
func TestArmoryOfficerNel(t *testing.T) {
	var nel, upgrade, topDeck ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(ct.Bind(&nel, ArmoryOfficerNel)),
			Hand: ct.Cards(
				ct.Bind(&upgrade, ct.Upgrade(ct.OfHouse(card.House.StarAlliance))),
			),
			Deck: ct.Cards(ct.Bind(&topDeck, ct.Creature())),
		},
	})

	h.P1.Play(upgrade)

	h.Expect(topDeck).At(ct.Hand)
}
