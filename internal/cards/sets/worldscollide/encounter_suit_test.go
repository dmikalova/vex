package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Encounter Suit
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//
//	This Creature gains, "After a Tactic is played but before it resolves, ward this Creature."
func TestEncounterSuit(t *testing.T) {
	t.Run("wards its host after an action card is played", func(t *testing.T) {
		var host ct.Card
		action := ct.Tactic(ct.OfHouse(card.House.StarAlliance))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						EncounterSuit,
					),
				),
				Hand: ct.Cards(action),
			},
		})

		h.P1.Play(action)

		if !h.Game().Warded(host.ID()) {
			t.Error("the host should be warded after an action card is played")
		}
	})

	t.Run("does not ward its host when a creature is played", func(t *testing.T) {
		var host ct.Card
		filler := ct.Creature(ct.OfHouse(card.House.StarAlliance))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						EncounterSuit,
					),
				),
				Hand: ct.Cards(filler),
			},
		})

		h.P1.Play(filler)

		if h.Game().Warded(host.ID()) {
			t.Error("playing a creature should not ward the host")
		}
	})
}
