package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bonesaw
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	If a friendly creature has been destroyed this turn, Bonesaw enters play ready.
func TestBonesaw(t *testing.T) {
	t.Run(
		"enters play ready after a friendly creature was destroyed this turn",
		func(t *testing.T) {
			var bonesaw, chump, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					Hand:  ct.Cards(ct.Bind(&bonesaw, Bonesaw)),
					InPlay: ct.Cards(
						ct.Bind(&chump, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(1))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(5)))),
				},
			})

			h.P1.Fight(chump, enemy)
			h.Expect(chump).At(ct.Discard)

			h.P1.Play(Bonesaw)
			if bonesaw.Exhausted() {
				t.Error("Bonesaw should enter play ready after a friendly creature died")
			}
		},
	)

	t.Run("enters play exhausted with no friendly deaths", func(t *testing.T) {
		var bonesaw ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&bonesaw, Bonesaw)),
			},
		})

		h.P1.Play(Bonesaw)
		if !bonesaw.Exhausted() {
			t.Error("Bonesaw should enter play exhausted with no friendly deaths")
		}
	})
}
