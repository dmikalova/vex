package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Subject Kirby
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Play/Fight/Reap: Play a non-Star Alliance creature.
func TestSubjectKirby(t *testing.T) {
	t.Run("reaping plays one off-house creature immediately", func(t *testing.T) {
		var kirby ct.Card
		marsCreature := ct.Creature(ct.OfHouse(card.House.Mars))
		marsArtifact := ct.Artifact(ct.OfHouse(card.House.Mars))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&kirby, SubjectKirby)),
				Hand:   ct.Cards(marsCreature, marsArtifact),
			},
		})

		// A non-Star Alliance card cannot be played on its own before reaping.
		h.P1.ExpectCannotPlay(marsCreature)

		// Reaping plays a creature at once: only the creature qualifies (the
		// artifact is barred), so it is played automatically and enters play.
		h.P1.Reap(kirby)

		h.Expect(marsCreature).At(ct.PlayArea)
		h.Expect(marsArtifact).At(ct.Hand)
	})
}
