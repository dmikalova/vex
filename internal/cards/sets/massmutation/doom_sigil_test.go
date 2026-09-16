package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Doom Sigil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	If there are no creatures in play, destroy Doom Sigil.
//	Each creature gains poison.
func TestDoomSigil(t *testing.T) {
	t.Run(
		"grants poison to every creature and self-destroys once the board empties",
		func(t *testing.T) {
			var sigil, friend, foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Shadows,
					InPlay: ct.Cards(
						ct.Bind(&sigil, DoomSigil),
						ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
					),
				},
				P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
			})

			if !h.Game().HasKeyword(friend.ID(), card.Keyword.Poison) {
				t.Error("friendly creature should gain poison from Doom Sigil")
			}
			if !h.Game().HasKeyword(foe.ID(), card.Keyword.Poison) {
				t.Error("enemy creature should gain poison from Doom Sigil")
			}
			h.Expect(sigil).At(ct.PlayArea)

			// Both 3-power creatures trade in a fight; poison ensures the foe dies. With
			// no creatures left in play, Doom Sigil destroys itself.
			h.P1.Fight(friend, foe)
			h.Expect(friend).At(ct.Discard)
			h.Expect(foe).At(ct.Discard)
			h.Expect(sigil).At(ct.Discard)
		},
	)

	t.Run("survives while a creature remains in play", func(t *testing.T) {
		var sigil ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&sigil, DoomSigil), ct.Creature()),
			},
			P2: ct.Side{},
		})

		h.Expect(sigil).At(ct.PlayArea)
	})
}
