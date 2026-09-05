package engine

import "testing"

// An abduction puts an enemy creature into the abductor's archives. Your archives
// are one of the three zones that may hold an enemy card, so the card needs no
// rider: leaving them, it goes to its owner's matching zone.

func TestAbductionText(t *testing.T) {
	single := PutFromPlay{
		Target:      Target{Kind: TargetChosenEnemyCreature},
		Destination: ToArchives.Yours(),
	}
	if got := single.Text(); got != "put an enemy creature into your archives" {
		t.Errorf("text = %q", got)
	}

	many := PutChosen{
		Amount:      3,
		UpTo:        true,
		Target:      Target{Kind: TargetEachEnemyCreature},
		Destination: ToArchives.Yours(),
	}
	if got := many.Text(); got != "put up to 3 enemy creatures into your archives" {
		t.Errorf("plural text = %q", got)
	}
}

func TestAbductionResolve(t *testing.T) {
	abduct := PutFromPlay{
		Target:      Target{Kind: TargetChosenEnemyCreature},
		Destination: ToArchives.Yours(),
	}

	// setup abducts player 1's sole creature into player 0's archives.
	setup := func(t *testing.T) (*Game, LocalID) {
		t.Helper()
		g := NewGame("A", "B", 1)
		prey := g.AddToBattleline(testCreature("prey", 3), 1)
		abduct.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if !g.State.Archives[0].contains(prey) {
			t.Fatalf(
				"prey should sit in the abductor's archives, got %v",
				g.State.Archives[0].slice(),
			)
		}
		return g, prey
	}

	t.Run("taken into hand", func(t *testing.T) {
		g, prey := setup(t)
		g.offerArchives(0) // the default chooser answers "Yes"
		if !g.State.Hand[1].contains(prey) {
			t.Error("prey should go to its owner's hand, not the abductor's")
		}
		if g.State.Hand[0].contains(prey) {
			t.Error("prey should not enter the abductor's hand")
		}
	})

	t.Run("archives discarded", func(t *testing.T) {
		g, prey := setup(t)
		g.discardArchives(0)
		if !g.State.Discard[1].contains(prey) {
			t.Error("a discarded abductee should go to its owner's discard pile")
		}
		if g.State.Discard[0].contains(prey) {
			t.Error("prey should not enter the abductor's discard pile")
		}
	})

	t.Run("ordinary archived card is untouched", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		own := g.AddToHand(testCreature("own", 3), 0)
		g.archiveFromHand(0, own)
		g.offerArchives(0)
		if !g.State.Hand[0].contains(own) {
			t.Error("a card archived normally returns to its own controller's hand")
		}
	})
}
