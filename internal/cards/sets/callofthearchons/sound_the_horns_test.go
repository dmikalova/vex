package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sound the Horns
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Discard cards from the top of your deck until you discard a Brobnar creature or run out of cards -> put it into your hand.
func TestSoundTheHorns(t *testing.T) {
	t.Run("digs to the first Brobnar creature and takes it", func(t *testing.T) {
		var horns, skipped, brute, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&horns, SoundTheHorns)),
				Deck: ct.Cards(
					ct.Bind(&skipped, ct.Tactic(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&brute, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
		})

		h.P1.Play(horns)

		h.Expect(brute).At(ct.Hand)
		h.Expect(skipped).At(ct.Discard)
		h.Expect(buried).At(ct.Deck)
	})

	t.Run("empties the deck when it finds nothing", func(t *testing.T) {
		var horns, offHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&horns, SoundTheHorns)),
				Deck: ct.Cards(
					ct.Bind(&offHouse, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Play(horns)

		h.Expect(offHouse).At(ct.Discard)
	})
}
