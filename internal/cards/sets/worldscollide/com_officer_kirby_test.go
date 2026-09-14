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
//	Play/Fight/Reap: Play a non-Star Alliance Artifact, Upgrade, or Tactic.
func TestComOfficerKirby(t *testing.T) {
	t.Run("reaping plays one off-house non-creature immediately", func(t *testing.T) {
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

		// A non-Star Alliance card cannot be played on its own before reaping.
		h.P1.ExpectCannotPlay(marsArtifact)

		// Reaping plays a non-creature at once: only the artifact qualifies (the
		// creature is barred), so it is played automatically and enters play.
		h.P1.Reap(kirby)

		h.Expect(marsArtifact).At(ct.PlayArea)
		h.Expect(marsCreature).At(ct.Hand)
	})
}
