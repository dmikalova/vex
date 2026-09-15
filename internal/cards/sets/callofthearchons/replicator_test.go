package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Replicator
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Mutant
//
//	Reap: Trigger the reap effect of another creature.
func TestReplicator(t *testing.T) {
	t.Run("reaps and fires another creature's reap effect", func(t *testing.T) {
		var replicator, faerie ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&replicator, Replicator),
					ct.Bind(&faerie, DewFaerie),
				),
			},
		})

		// Dew Faerie is the only other creature carrying a Reap ability, so it is
		// chosen without a prompt. It gains its Æmber without exhausting.
		h.P1.Reap(replicator)

		h.P1.ExpectAmber(2) // 1 for Replicator's reap + 1 from Dew Faerie
		h.Expect(faerie).Ready()
	})
}
