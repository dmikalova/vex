package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Com. Officer Kirby
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: You may play a non-Star Alliance artifact, upgrade, or Tactic this turn.
func TestComOfficerKirby(t *testing.T) {
	t.Run("reaping frees one off-house non-creature this turn", func(t *testing.T) {
		var kirby ct.Card
		marsArtifact := ct.Artifact(ct.OfHouse(card.House.Mars))
		marsCreature := ct.Creature(ct.OfHouse(card.House.Mars))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&kirby, ComOfficerKirby)),
				Hand:   ct.Cards(marsArtifact, marsCreature),
			},
		})

		// A non-Star Alliance card cannot be played before the grant.
		h.P1.ExpectCannotPlay(marsArtifact)

		h.P1.Reap(kirby)

		// The freed non-creature may now be played; the creature stays barred.
		h.P1.ExpectCannotPlay(marsCreature)
		h.P1.Play(marsArtifact)
	})
}
