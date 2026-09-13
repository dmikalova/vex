package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Transporter Platform
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	Action: Put a friendly Creature and each Upgrade attached to it into its owner's hand.
func TestTransporterPlatform(t *testing.T) {
	var creature, upgrade ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.StarAlliance,
			InPlay: ct.Cards(
				TransporterPlatform,
				ct.Bind(&creature, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3))),
			),
		},
	})
	upgrade = creature.Attach(ct.Upgrade())

	h.P1.UseAction(TransporterPlatform) // the sole friendly creature is auto-chosen

	h.Expect(creature).At(ct.Hand)
	h.Expect(upgrade).At(ct.Hand)
}
