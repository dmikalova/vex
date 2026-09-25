package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Whisper
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Lose 1 Æmber -> destroy a creature.
func TestWhisper(t *testing.T) {
	t.Run("loses 1 Æmber to destroy a creature", func(t *testing.T) {
		var whisper, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Amber:  1,
				InPlay: ct.Cards(ct.Bind(&whisper, Whisper)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5))))},
		})

		h.P1.UseAction(whisper)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.Discard)
		h.P1.ExpectAmber(0)
	})
}
